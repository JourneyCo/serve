import {AfterViewInit, Component, OnInit, ViewChild} from '@angular/core';
import { CommonModule } from "@angular/common";
import { RouterModule, Router } from "@angular/router";
import { GoogleMapsModule } from "@angular/google-maps";
import {HelperService} from '@services';
import { Project } from '@models';
import { ProjectService } from '@services';
import { Marker } from '@models';
import {Subject} from 'rxjs';
import {MaterialModule} from '@material';
import {MatPaginator, PageEvent} from '@angular/material/paginator';
import {MatTableDataSource} from '@angular/material/table';
import {MatSort} from '@angular/material/sort';

@Component({
  selector: "app-project-list",
  standalone: true,
  imports: [
    CommonModule,
    RouterModule,
    GoogleMapsModule,
    MaterialModule
  ],
  templateUrl: "./projects.component.html",
  styleUrls: ["./projects.component.scss"],
})
export class ProjectsComponent implements OnInit, AfterViewInit {
  displayedColumns: string[] = [
    "title",
    "projectLocation",
    "time",
    "capacity",
    "actions",
  ];
  dataSource = new MatTableDataSource<Project>([]);
  isLoading = true;
  eventsSubject: Subject<any> = new Subject<any>();
  pageIndex: number = 0;
  pageSize: number = 10;


  @ViewChild(MatPaginator) paginator!: MatPaginator;
  @ViewChild(MatSort) sort!: MatSort;

  // Google Maps properties
  mapOptions: google.maps.MapOptions = {
    center: { lat: 40.0, lng: -95.0 },
    zoom: 12,
  };
  markers: Marker[] = [];
  mapsLoaded = false;

  constructor(
    private projectService: ProjectService,
    private router: Router,
    private helper: HelperService,
  ) {}

  ngOnInit(): void {
    // Google Maps API is automatically loaded by the Angular Google Maps module
    this.mapsLoaded = true;
    this.loadProjects();

  }

  ngAfterViewInit() {
    setTimeout(() => {
      this.dataSource.paginator = this.paginator;
      this.dataSource.sort = this.sort;
    });
  }

  loadProjects(): void {
    this.isLoading = true;
    this.projectService.getProjects().subscribe(
      (projects) => {
        this.dataSource.data = projects;

        // Create markers for projects with valid coordinates
        this.markers = [];
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
          this.markers = validProjects.map((project, index) => ({
            position: {
              lat: project.latitude || 0,
              lng: project.longitude || 0,
            },
            label: String(index + 1),
            title: project.title,
          }));
        }

        const startIndex = this.pageIndex * this.pageSize;
        const endIndex = startIndex + this.pageSize;
        this.dataSource.data = Array.from({ length: this.dataSource.data.length }, (_, i) =>
          this.dataSource.data[i]).slice(startIndex, endIndex);

        this.isLoading = false;
      },
      (error) => {
        console.error("Error loading projects:", error);
        this.isLoading = false;
        this.helper.showError(
          "Error loading projects. Please try again.");
      },
    );
  }

  applyFilter(event: Event): void {
    const filterValue = (event.target as HTMLInputElement).value;
    this.dataSource.filter = filterValue.trim().toLowerCase();

    if (this.dataSource.paginator) {
      this.dataSource.paginator.firstPage();
    }
  }

  viewProject(project: Project): void {
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

  formatDate(date: string): string {
    return new Date(date).toLocaleDateString();
  }

  formatTime(time: string): string {
    if (!time) return "";

    // Assuming time is in format "HH:MM:SS" or "HH:MM"
    const parts = time.split(":");
    if (parts.length < 2) return time;

    // Convert to 12-hour format
    let hours = parseInt(parts[0], 10);
    const minutes = parts[1];
    const ampm = hours >= 12 ? "PM" : "AM";

    hours = hours % 12;
    hours = hours ? hours : 12; // the hour '0' should be '12'

    return `${hours}:${minutes} ${ampm}`;
  }

  markerClicked(marker: any, project: Project): void {
    this.viewProject(project);
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

  handlePageEvent(event: PageEvent) {
    this.pageSize = event.pageSize;
    this.pageIndex = event.pageIndex;
    this.loadProjects();
  }

  trackByFn(index: number, project: Project): number {
    return project.id;
  }
}
