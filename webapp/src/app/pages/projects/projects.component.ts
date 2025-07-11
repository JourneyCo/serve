import { AfterViewChecked, AfterViewInit, Component, OnInit, ViewChild } from '@angular/core';
import { CommonModule } from '@angular/common';
import { Router, RouterModule } from '@angular/router';
import { GoogleMapsModule, MapAdvancedMarker, MapInfoWindow } from '@angular/google-maps';
import { HelperService, ProjectService } from '@services';
import { Project } from '@models';
import { MaterialModule } from '@material';
import { MatTableDataSource } from '@angular/material/table';
import { MatSort } from '@angular/material/sort';

@Component({
  selector: "app-project-list",
  standalone: true,
  imports: [
    CommonModule,
    RouterModule,
    GoogleMapsModule,
    MaterialModule,
  ],
  templateUrl: "./projects.component.html",
  styleUrls: ["./projects.component.scss"],
})
export class ProjectsComponent implements OnInit, AfterViewInit, AfterViewChecked {
  displayedColumns: string[] = [
    "title",
    "type",
    "area",
    "time",
    "capacity",
  ];
  dataSource = new MatTableDataSource<Project>([]);
  isLoading = true;


  @ViewChild(MatSort) sort!: MatSort;

  // Google Maps properties
  mapOptions: google.maps.MapOptions = {
    center: { lat: 40.0, lng: -95.0 },
    zoom: 12,
  };
  markers: any[];
  mapsLoaded = false;
  @ViewChild(MapInfoWindow) infoWindow: MapInfoWindow;
  selectedProject: Project | null;
  types: Record<number, string> = {};

  constructor(
    private projectService: ProjectService,
    private router: Router,
    private helper: HelperService,
  ) { }

  ngOnInit(): void {
    this.mapsLoaded = true;
    this.loadTypes();
    this.loadProjects();
  }

  loadTypes(): void {
    this.projectService.getTypes().subscribe({
      next: (typesData) => {
        // Convert array to Record<number, string> format
        this.types = {};
        typesData.forEach(type => {
          this.types[type.id] = type.name;
        });
      },
      error: (error: any) => {
        console.error("Error loading types:", error);
        this.helper.showError("Error loading project types");
      }
    });
  }

  ngAfterViewInit() {
    setTimeout(() => {
      this.dataSource.sort = this.sort;
    });
  }

  private sortInitialized = false;

  ngAfterViewChecked() {
    if (!this.sortInitialized && this.sort && !this.isLoading) {
      this.dataSource.sort = this.sort;
      this.sortInitialized = true;
    }
  }

  loadProjects(): void {
    this.isLoading = true;
    this.projectService.getProjects().subscribe({
      next: (data) => {
        let projects = this.sortProjects(data);


        this.dataSource.data = projects;

        // Create markers for projects with valid coordinates
        const validProjects = projects.filter((p) => p.latitude && p.longitude);

        if (validProjects.length > 0) {
          // Get average lat/lng for centering the google-map
          const totalLat = validProjects.reduce(
            (sum, p) => sum + (p.latitude || 0),
            0,
          );
          const totalLng = validProjects.reduce(
            (sum, p) => sum + (p.longitude || 0),
            0,
          );
          const avgLat = totalLat / validProjects.length;
          const avgLng = totalLng / validProjects.length;

          // Update google-map center
          this.mapOptions = {
            ...this.mapOptions,
            center: { lat: avgLat, lng: avgLng },
            zoom: validProjects.length > 1 ? 10 : 12, // Zoom level based on number of markers
          };

          // Create markers
          this.markers = [];
          const markerMap = new Map<string, number>();
          const OFFSET = 0.0001; // Adjust this value as needed

          validProjects.forEach((project) => {
            const key = `${project.latitude},${project.longitude}`;
            let offsetIndex = 0;
            if (markerMap.has(key)) {
              offsetIndex = markerMap.get(key)! + 1;
              markerMap.set(key, offsetIndex);
            } else {
              markerMap.set(key, 0);
            }
            // Offset calculation: spread markers in a spiral for better visibility

            // Calculate small offsets in a circle pattern to prevent overlap
            let latOffset = 0;
            let lngOffset = 0;

            if (offsetIndex > 0) {
              const angle = (offsetIndex - 1) * (360 / 8); // Distribute up to 8 markers in a circle
              const radian = (angle * Math.PI) / 180;
              latOffset = Math.cos(radian) * OFFSET;
              lngOffset = Math.sin(radian) * OFFSET;
            }

            // @ts-ignore
            const mark: MapAdvancedMarker = {
              position: {
                lat: (project.latitude || 0) + latOffset,
                lng: (project.longitude || 0) + lngOffset,
              },
              title: project.title,
            }
            const set = {
              mark: mark,
              project: project
            }
            this.markers.push(set);
          })
        }

        this.isLoading = false;
      },
      error: (error) => {
        console.error("Error loading projects:", error);
        this.isLoading = false;
        this.helper.showError(
          "Error loading projects. Please try again.");
      },
    });
  }

  applyFilter(event: Event): void {
    const filterValue = (event.target as HTMLInputElement).value;
    this.dataSource.filter = filterValue.trim().toLowerCase();
  }

  viewProject(project: Project | null): void {
    if (!project) {
      return
    }
    this.router.navigate(["/projects", project.id]);
  }

  getCapacityPercentage(project: Project): number {
    return (project.current_registrations / project.max_capacity) * 100;
  }

  getCapacityStatusClass(project: Project): string {
    const percentage = this.getCapacityPercentage(project);

    if (percentage >= 100) {
      return "full";
    } else if (percentage >= 80) {
      return "almost-full";
    } else if (percentage >= 50) {
      return "filling";
    } else {
      return "available";
    }
  }

  markerClicked(marker: MapAdvancedMarker, mark: any): void {
    this.selectedProject = mark.project
    this.infoWindow.open(marker);
  }

  centerMap(): void {
    if (this.markers.length > 0) {
      // Recalculate center and zoom to fit all markers
      const bounds = new google.maps.LatLngBounds();
      this.markers.forEach((marker) => {
        bounds.extend(marker.position);
      });

      this.mapOptions = {
        ...this.mapOptions,
        center: {
          lat: bounds.getCenter().lat(),
          lng: bounds.getCenter().lng(),
        },
      };
    }
  }
  trackByFn(index: number, project: Project): number {
    return project.id;
  }

  sortProjects(p_list: Project[]) {
    return p_list.sort((a, b) => {
      // 1. Status: 'open' comes first
      const aOpen = a.status === 'open';
      const bOpen = b.status === 'open';
      if (aOpen !== bOpen) {
        return aOpen ? -1 : 1;
      }

      // 2. Fullness: not-full comes before full
      const aFull = a.current_registrations >= a.max_capacity;
      const bFull = b.current_registrations >= b.max_capacity;
      if (aFull !== bFull) {
        return aFull ? 1 : -1;
      }

      // 3. Alphabetical by title
      return a.title.localeCompare(b.title);
    });
  }
}