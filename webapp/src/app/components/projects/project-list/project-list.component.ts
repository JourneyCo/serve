import { Component, OnInit, ViewChild } from "@angular/core";
import { CommonModule } from "@angular/common";
import { RouterModule, Router } from "@angular/router";
import { MatPaginator, MatPaginatorModule } from "@angular/material/paginator";
import { MatSort, MatSortModule } from "@angular/material/sort";
import { MatTableModule, MatTableDataSource } from "@angular/material/table";
import { MatCardModule } from "@angular/material/card";
import { MatFormFieldModule } from "@angular/material/form-field";
import { MatInputModule } from "@angular/material/input";
import { MatButtonModule } from "@angular/material/button";
import { MatIconModule } from "@angular/material/icon";
import { MatProgressSpinnerModule } from "@angular/material/progress-spinner";
import { MatSnackBar, MatSnackBarModule } from "@angular/material/snack-bar";
import { MatButtonToggleModule } from "@angular/material/button-toggle";
import { GoogleMapsModule } from "@angular/google-maps";
import { GoogleMapsApiService } from '@services';
import { Project } from '@models';
import { ProjectService } from '@services';
import { GoogleMapComponent } from "../../google-map/google-map.component";
import { Marker } from '@models';
import { Subject } from "rxjs";

@Component({
  selector: "app-project-list",
  standalone: true,
  imports: [
    CommonModule,
    RouterModule,
    MatTableModule,
    MatSortModule,
    MatPaginatorModule,
    MatCardModule,
    MatFormFieldModule,
    MatInputModule,
    MatButtonModule,
    MatIconModule,
    MatProgressSpinnerModule,
    MatSnackBarModule,
    MatButtonToggleModule,
    GoogleMapsModule,
  ],
  templateUrl: "./project-list.component.html",
  styleUrls: ["./project-list.component.scss"],
})
export class ProjectListComponent implements OnInit {
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
    private snackBar: MatSnackBar,
    private mapsApiService: GoogleMapsApiService,
  ) {}

  ngOnInit(): void {
    // Google Maps API is automatically loaded by the Angular Google Maps module
    this.mapsLoaded = true;
    this.loadProjects();
  }

  ngAfterViewInit() {
    this.dataSource.paginator = this.paginator;
    this.dataSource.sort = this.sort;
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

        this.isLoading = false;
      },
      (error) => {
        console.error("Error loading projects:", error);
        this.isLoading = false;
        this.snackBar.open(
          "Error loading projects. Please try again.",
          "Close",
          {
            duration: 5000,
          },
        );
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

  isAtCapacity(project: Project): boolean {
    return project.current_registrations >= project.max_capacity;
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

  trackByFn(index: number, project: Project): number {
    return project.id;
  }
}
