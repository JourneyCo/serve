import {Component, OnInit, ViewChild, TemplateRef} from "@angular/core";
import {CommonModule} from "@angular/common";
import { FormsModule } from "@angular/forms";
import { ActivatedRoute, Router, RouterModule } from "@angular/router";
import { MatDialog} from '@angular/material/dialog';
import { GoogleMapsModule } from "@angular/google-maps";
import { AuthService, ProjectService, HelperService, RegistrationService } from "@services";
import {Observable, forkJoin, of, Subscription} from 'rxjs';
import {User, Project, Ages, Categories, Supplies, Tools, Skills} from '@models';
import {AdminProjectPanelComponent, RegistrationDialogComponent} from '@components';
import { MaterialModule } from '@material';

@Component({
  selector: "app-project-detail",
  standalone: true,
  imports: [
    CommonModule,
    FormsModule,
    RouterModule,
    GoogleMapsModule,
    AdminProjectPanelComponent,
    MaterialModule
  ],
  templateUrl: "./project-detail.component.html",
  styleUrls: ["./project-detail.component.scss"],
})

export class ProjectDetailComponent implements OnInit {
  project: Project | null = null;
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
  skills = Skills
  serve_date: Date;
  registrationSubscription: Subscription;

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
    private dialog: MatDialog,
    private helper: HelperService,
    private registrationService: RegistrationService,
  ) {
    this.serve_date = helper.GetServeDate();
    this.isAdmin = authService.isAdmin();
  }

  ngOnInit(): void {
    // Google Maps API is automatically loaded by the Angular Google Maps module
    // Just load the project data directly
    this.loadProjectData();
    this.registrationSubscription = this.registrationService.registrationChange$.subscribe(() => {
      this.loadProjectData();
    });
  }

  loadProjectData(): void {
    this.isLoading = true;

    // Get the project ID from the route
    const project_id = this.route.snapshot.paramMap.get("id");
    if (!project_id || isNaN(+project_id)) {
      this.helper.showError("Invalid project ID");
      this.router.navigate(["/projects"]);
      return;
    }

    // Get current user and check if admin
    this.authService.getCurrentUser().subscribe({
      next: (user) => {
        this.currentUser = user;

        // Load project data
        this.loadProjectDetails(+project_id);
      },
      error: (error: any) => {
        console.error("Error getting user:", error);
        this.helper.showError("Error loading user information");
        this.isLoading = false;
      },
    });
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
    }).subscribe({
      next: (result) => {
        this.project = result.project;

        // Check if user is registered for this project
        this.isRegistered = result.userRegs?.some(
          (reg) => reg.project_id === project_id && reg.status === "registered",
        );

        // Set up google-map marker if coordinates are available
        this.updateMapMarker();

        this.isLoading = false;
      },
      error: (error: any) => {
        console.error("Error loading project details:", error);
        this.helper.showError("Error loading project details");
        this.isLoading = false;
      },
    });
  }

  openRegistrationForm(): void {
    // Open dialog
    this.dialogRef = this.dialog.open(RegistrationDialogComponent, {
      width: "500px",
      disableClose: false,
      data: { project: this.project, user: this.currentUser }
    });

    // Handle dialog closure
    this.dialogRef.afterClosed().subscribe((data: any) => {
      // If dialog was dismissed, do nothing
      if (!data.success) {
        this.dialogRef = null;
        return
      }

      this.registerForProject(data.values)
    });
  }

  registerForProject(data: any): void {
    if (!this.project) return;

    this.loadingRegistration = true;
    this.registrationError = "";

    // Close dialog if open
    if (this.dialogRef) {
      this.dialogRef.close(true);
    }

    const body = {
      guest_count: data.guest_count,
      email: data.email,
      phone: data.phone,
      first_name: data.first_name,
      last_name: data.last_name,
      lead_interest: data.lead_interest,
    };

    this.projectService.registerForProject(this.project.id, body).subscribe({
      next: () => {
        this.loadingRegistration = false;
        this.isRegistered = true;
        this.helper.showSuccess("Successfully registered for the project");

        // Reload project to get updated capacity
        this.loadProjectDetails(this.project!.id);
        this.registrationService.triggerRegistrationChange();
      },
      error: (error: any) => {
        this.loadingRegistration = false;
        console.error("Registration error:", error);
        this.registrationError =
          error.error?.error || "Failed to register for the project";
        this.helper.showError(this.registrationError);
      },
    });
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

    this.projectService.cancelRegistration(this.project.id).subscribe({
      next: () => {
        this.loadingRegistration = false;
        this.isRegistered = false;
        this.helper.showSuccess("Registration cancelled successfully");

        // Reload project to get updated capacity
        this.loadProjectDetails(this.project!.id);
        this.registrationService.triggerRegistrationChange();
      },
      error: (error: any) => {
        this.loadingRegistration = false;
        console.error("Cancellation error:", error);
        this.helper.showError(error.error?.error || "Failed to cancel registration");
      },
    });
  }

  getDaysUntilStart(): number {
    if (!this.project) return 0;

    const today = new Date();
    today.setHours(0, 0, 0, 0);

    const diffTime = this.helper.GetServeDate().getTime() - today.getTime();
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

  updateMapMarker(): void {
    if (!this.project || !this.project.latitude || !this.project.longitude) {
      this.markerPosition = null;
      return
    }

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
    }
}
