import { Component, OnInit, ViewChild, TemplateRef } from "@angular/core";
import { CommonModule } from "@angular/common";
import { FormsModule } from "@angular/forms";
import { ActivatedRoute, Router, RouterModule } from "@angular/router";
import { MatSnackBar, MatSnackBarModule } from "@angular/material/snack-bar";
import { MatCardModule } from "@angular/material/card";
import { MatButtonModule } from "@angular/material/button";
import { MatIconModule } from "@angular/material/icon";
import { MatProgressSpinnerModule } from "@angular/material/progress-spinner";
import { MatListModule } from "@angular/material/list";
import { MatDividerModule } from "@angular/material/divider";
import { MatChipsModule } from "@angular/material/chips";
import { MatBadgeModule } from "@angular/material/badge";
import { MatProgressBarModule } from "@angular/material/progress-bar";
import { MatFormFieldModule } from "@angular/material/form-field";
import { MatInputModule } from "@angular/material/input";
import { MatCheckboxModule } from "@angular/material/checkbox";
import { MatDialogModule, MatDialog } from "@angular/material/dialog";
import { GoogleMapsModule } from "@angular/google-maps";
import { ProjectService } from "../../../services/project.service";
import { AuthService } from "../../../services/auth.service";
import { GoogleMapsApiService } from "../../../services/google-maps-api.service";
import { Project } from "../../../models/project.model";
import { Registration } from "../../../models/registration.model";
import { User } from "../../../models/user.model";
import { Observable, forkJoin, of } from "rxjs";
import { tap } from "rxjs/operators";
import { catchError, map, switchMap } from "rxjs/operators";
import {MatTable, MatTableDataSource} from "@angular/material/table";

@Component({
  selector: "app-project-detail",
  standalone: true,
  imports: [
    CommonModule,
    FormsModule,
    RouterModule,
    MatCardModule,
    MatButtonModule,
    MatIconModule,
    MatProgressSpinnerModule,
    MatListModule,
    MatDividerModule,
    MatChipsModule,
    MatBadgeModule,
    MatProgressBarModule,
    MatSnackBarModule,
    MatFormFieldModule,
    MatInputModule,
    MatCheckboxModule,
    MatDialogModule,
    GoogleMapsModule,
    MatTable,
  ],
  templateUrl: "./project-detail.component.html",
  styleUrls: ["./project-detail.component.scss"],
})
export class ProjectDetailComponent implements OnInit {
  registrationsColumns = ['avatar', 'name', 'email', 'status', 'registrationDate'];
  project: Project | null = null;
  registrations: Registration[] = [];
  registrationsDataSource = new MatTableDataSource<Registration>();
  currentUser: User | null = null;
  isAdmin = false;
  isRegistered = false;
  isLoading = true;
  loadingRegistration = false;
  registrationError = "";

  // Registration form properties
  guestCount: number = 0;
  isProjectLead: boolean = false;
  firstName: string = '';
  lastName: string = '';
  phone: string = '';
  contactEmail: string = '';

  // Google Maps properties
  mapOptions: google.maps.MapOptions = {
    zoom: 14,
    mapTypeControl: false,
    streetViewControl: false,
    fullscreenControl: true,
  };
  markerOptions: any = {
    draggable: false,
  };
  markerPosition: google.maps.LatLngLiteral | null = null;

  @ViewChild("registrationDialog")
  registrationDialogTemplate!: TemplateRef<any>;
  dialogRef: any;

  constructor(
    private route: ActivatedRoute,
    private router: Router,
    private projectService: ProjectService,
    private authService: AuthService,
    private mapsApiService: GoogleMapsApiService,
    private snackBar: MatSnackBar,
    private dialog: MatDialog,
  ) {}

  ngOnInit(): void {
    // Google Maps API is automatically loaded by the Angular Google Maps module
    // Just load the project data directly
    this.loadProjectData();
  }

  loadProjectData(): void {
    this.isLoading = true;

    // Get the project ID from the route
    const projectId = this.route.snapshot.paramMap.get("id");
    if (!projectId || isNaN(+projectId)) {
      this.showError("Invalid project ID");
      this.router.navigate(["/projects"]);
      return;
    }

    // Get current user and check if admin
    this.authService.getCurrentUser().subscribe(
      (user) => {
        this.currentUser = user;
        // this.isAdmin = user?.isAdmin || false;
        this.isAdmin = true;


        // Load project data
        this.loadProjectDetails(+projectId);
      },
      (error: any) => {
        console.error("Error getting user:", error);
        this.showError("Error loading user information");
        this.isLoading = false;
      },
    );
  }

  loadProjectDetails(projectId: number): void {
    // Create observables for project and user registrations
    const projectObs = this.projectService.getProject(projectId);
    const registrationsObs = this.projectService.getUserRegistrations();
    const projectRegistrationsObs = this.isAdmin
      ? this.projectService.getProjectRegistrations(projectId)
      : of([]);

    // Use forkJoin to get all data at once
    forkJoin({
      project: projectObs,
      userRegs: registrationsObs,
      projectRegs: projectRegistrationsObs,
    }).subscribe(
      (result) => {
        this.project = result.project;

        // Check if user is registered for this project
        this.isRegistered = result.userRegs?.some(
          (reg) => reg.projectId === projectId && reg.status === "registered",
        );

        // If admin, get all registrations for the project
        if (this.isAdmin) {
          this.registrations = result.projectRegs;
          this.registrationsDataSource.data = this.registrations;
        }

        // Set up map marker if coordinates are available
        this.updateMapMarker();

        this.isLoading = false;
      },
      (error: any) => {
        console.error("Error loading project details:", error);
        this.showError("Error loading project details");
        this.isLoading = false;
      },
    );
  }

  openRegistrationForm(): void {
    // Reset form values
    this.guestCount = 0;
    this.isProjectLead = false;

    // Open dialog
    this.dialogRef = this.dialog.open(this.registrationDialogTemplate, {
      width: "500px",
      disableClose: false,
    });

    // Handle dialog closure
    this.dialogRef.afterClosed().subscribe((result: boolean) => {
      // If dialog was dismissed, do nothing
      if (!result) {
        this.dialogRef = null;
      }
    });
  }

  registerForProject(): void {
    if (!this.project) return;

    this.loadingRegistration = true;
    this.registrationError = "";

    // Close dialog if open
    if (this.dialogRef) {
      this.dialogRef.close(true);
    }

    this.projectService
      .registerForProject(this.project.id, this.guestCount, this.isProjectLead)
      .subscribe(
        (_) => {
          this.loadingRegistration = false;
          this.isRegistered = true;
          this.showSuccess("Successfully registered for the project");

          // Reload project to get updated capacity
          this.loadProjectDetails(this.project!.id);
        },
        (error: any) => {
          this.loadingRegistration = false;
          console.error("Registration error:", error);
          this.registrationError =
            error.error?.error || "Failed to register for the project";
          this.showError(this.registrationError);
        },
      );
  }

  cancelRegistration(): void {
    if (!this.project) return;

    this.loadingRegistration = true;

    this.projectService.cancelRegistration(this.project.id).subscribe(
      (_) => {
        this.loadingRegistration = false;
        this.isRegistered = false;
        this.showSuccess("Registration cancelled successfully");

        // Reload project to get updated capacity
        this.loadProjectDetails(this.project!.id);
      },
      (error: any) => {
        this.loadingRegistration = false;
        console.error("Cancellation error:", error);
        this.showError(error.error?.error || "Failed to cancel registration");
      },
    );
  }

  formatDate(date: string): string {
    return new Date(date).toLocaleDateString("en-US", {
      weekday: "long",
      year: "numeric",
      month: "long",
      day: "numeric",
    });
  }

  formatTime(time: string): string {
    // Convert HH:MM:SS to a formatted time string (e.g. "9:00 AM")
    const [hours, minutes] = time.split(":");
    const hour = parseInt(hours, 10);
    const minute = parseInt(minutes, 10);
    const ampm = hour >= 12 ? "PM" : "AM";
    const formattedHour = hour % 12 || 12; // Convert 0 to 12 for display
    return `${formattedHour}:${minute.toString().padStart(2, "0")} ${ampm}`;
  }

  getDaysUntilStart(): number {
    if (!this.project) return 0;

    const today = new Date();
    today.setHours(0, 0, 0, 0);

    const projectDate = new Date(this.project.projectDate);
    projectDate.setHours(0, 0, 0, 0);

    const diffTime = projectDate.getTime() - today.getTime();
    return Math.ceil(diffTime / (1000 * 60 * 60 * 24));
  }

  isProjectFull(): boolean {
    return this.project
      ? this.project.currentRegistrations >= this.project.maxCapacity
      : false;
  }

  getCapacityPercentage(): number {
    if (!this.project) return 0;
    return (this.project.currentRegistrations / this.project.maxCapacity) * 100;
  }

  private showSuccess(message: string): void {
    this.snackBar.open(message, "Close", {
      duration: 5000,
      panelClass: ["success-snackbar"],
    });
  }

  private showError(message: string): void {
    this.snackBar.open(message, "Close", {
      duration: 5000,
      panelClass: ["error-snackbar"],
    });
  }

  updateMapMarker(): void {
    if (this.project && this.project.latitude && this.project.longitude) {
      // Set the marker position and map center
      this.markerPosition = {
        lat: this.project.latitude,
        lng: this.project.longitude,
      };

      // Update map options to center on project location
      this.mapOptions = {
        ...this.mapOptions,
        center: this.markerPosition,
      };
    } else {
      this.markerPosition = null;
    }
  }
}
