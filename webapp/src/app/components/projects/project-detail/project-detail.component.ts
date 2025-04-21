import {Component, OnInit, ViewChild, TemplateRef} from "@angular/core";
import {CommonModule} from "@angular/common";
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
import { AuthService, ProjectService, HelperService } from "@services";
import { Observable, forkJoin, of } from "rxjs";
import {
  MatTable,
  MatTableDataSource,
  MatTableModule,
} from "@angular/material/table";
import {User, Registration, Project, Ages, Categories, Supplies, Tools} from "@models";

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
    MatTableModule,
  ],
  templateUrl: "./project-detail.component.html",
  styleUrls: ["./project-detail.component.scss"],
})
export class ProjectDetailComponent implements OnInit {
  registrationsColumns = ["name", "email", "phone", "registrationDate"];
  project: Project | null = null;
  registrations: Registration[] = [];
  registrationsDataSource = new MatTableDataSource<Registration>();
  currentUser: User | null = null;
  isAdmin: Observable<boolean> | undefined;
  isRegistered = false;
  isLoading = true;
  loadingRegistration = false;
  registrationError = "";
  tools = Tools
  supplies = Supplies
  categories = Categories
  ages = Ages
  serve_date: Date;

  // Registration form properties
  guest_count: number = 0;
  lead_interest: boolean = false;
  first_name: string = "";
  last_name: string = "";
  phone: string = "";
  email: string = "";

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
  @ViewChild("cancellationDialog")
  cancellationDialogTemplate!: TemplateRef<any>;
  dialogRef: any;

  constructor(
    private route: ActivatedRoute,
    private router: Router,
    private projectService: ProjectService,
    private authService: AuthService,
    private snackBar: MatSnackBar,
    private dialog: MatDialog,
    private helper: HelperService
  ) {
    this.serve_date = helper.GetServeDate();
  }

  ngOnInit(): void {
    // Google Maps API is automatically loaded by the Angular Google Maps module
    // Just load the project data directly
    this.loadProjectData();
  }

  loadProjectData(): void {
    this.isLoading = true;

    // Get the project ID from the route
    const project_id = this.route.snapshot.paramMap.get("id");
    if (!project_id || isNaN(+project_id)) {
      this.showError("Invalid project ID");
      this.router.navigate(["/projects"]);
      return;
    }

    // Get current user and check if admin
    this.authService.getCurrentUser().subscribe(
      (user) => {
        this.currentUser = user;

        // Load project data
        this.loadProjectDetails(+project_id);
      },
      (error: any) => {
        console.error("Error getting user:", error);
        this.showError("Error loading user information");
        this.isLoading = false;
      },
    );
  }

  loadProjectDetails(project_id: number): void {
    // Create observables for project and user registrations
    const projectObs = this.projectService.getProject(project_id);
    const registrationsObs = this.projectService.getUserRegistrations();
    const projectRegistrationsObs = this.isAdmin
      ? this.projectService.getProjectRegistrations(project_id)
      : of([]);

    // Use forkJoin to get all data at once
    forkJoin({
      project: projectObs,
      userRegs: registrationsObs,
      projectRegs: projectRegistrationsObs,
    }).subscribe(
      (result) => {
        this.project = result.project;
        console.log(this.project);

        // Check if user is registered for this project
        this.isRegistered = result.userRegs?.some(
          (reg) => reg.project_id === project_id && reg.status === "registered",
        );

        // If admin, get all registrations for the project
        if (this.isAdmin) {
          this.registrations = result.projectRegs;
          this.registrationsDataSource.data = this.registrations;
        }

        // Set up google-map marker if coordinates are available
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
    this.guest_count = 0;
    this.lead_interest = false;

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

    const body = {
      guest_count: this.guest_count,
      email: this.email,
      phone: this.phone,
      first_name: this.first_name,
      last_name: this.last_name,
      lead_interest: this.lead_interest,
    };

    this.projectService.registerForProject(this.project.id, body).subscribe(
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

  openCancellationDialog(): void {
    this.dialogRef = this.dialog.open(this.cancellationDialogTemplate, {
      width: "400px",
      disableClose: false
    });
  }

  confirmCancellation(): void {
    if (this.dialogRef) {
      this.dialogRef.close();
    }
    this.cancelRegistration();
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

  getDaysUntilStart(): number {
    if (!this.project) return 0;

    const today = new Date();
    today.setHours(0, 0, 0, 0);

    const diffTime = this.serve_date.getTime() - today.getTime();
    return Math.ceil(diffTime / (1000 * 60 * 60 * 24));
  }

  isProjectFull(): boolean {
    return this.project
      ? this.project.current_registrations >= this.project.max_capacity
      : false;
  }

  getCapacityPercentage(): number {
    if (!this.project) return 0;
    return (this.project.current_registrations / this.project.max_capacity) * 100;
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
      // Set the marker position and google-map center
      this.markerPosition = {
        lat: this.project.latitude,
        lng: this.project.longitude,
      };

      // Update google-map options to center on project location
      this.mapOptions = {
        ...this.mapOptions,
        center: this.markerPosition,
      };
    } else {
      this.markerPosition = null;
    }
  }
}
