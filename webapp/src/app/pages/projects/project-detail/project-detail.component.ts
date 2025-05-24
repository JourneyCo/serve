import {Component, OnInit, TemplateRef, ViewChild} from '@angular/core';
import {CommonModule} from '@angular/common';
import {FormsModule} from '@angular/forms';
import {ActivatedRoute, Router, RouterModule} from '@angular/router';
import {MatDialog} from '@angular/material/dialog';
import {GoogleMapsModule} from '@angular/google-maps';
import {AuthService, HelperService, ProjectService, RegistrationService} from '@services';
import {Observable, Subscription} from 'rxjs';
import {Ages, Categories, Project, Registration, User} from '@models';
import {
  AdminProjectPanelComponent,
  EditGuestCountDialogComponent,
  AlreadyRegisteredDialogComponent,
  AlreadyRegisteredElsewhereDialogComponent
} from '@components';
import {MaterialModule} from '@material';
import {NgxLinkifyjsModule, NgxLinkifyjsService} from 'ngx-linkifyjs-v2';


@Component({
  selector: "app-project-detail",
  standalone: true,
  imports: [
    CommonModule,
    FormsModule,
    RouterModule,
    GoogleMapsModule,
    AdminProjectPanelComponent,
    MaterialModule,
    NgxLinkifyjsModule,
  ],
  templateUrl: "./project-detail.component.html",
  styleUrls: ["./project-detail.component.scss"],
})

export class ProjectDetailComponent implements OnInit {
  project: Project | null = null;
  currentUser: User | null = null;
  isRegistered = false;
  isLoading = true;
  loadingRegistration = false;
  registrationError = "";
  categories = Categories
  ages = Ages
  serve_date: Date;
  registrationSubscription: Subscription;
  isAdmin: Observable<boolean>;
  userSignedIn: Observable<boolean>;

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
    private linkifyService: NgxLinkifyjsService,
  ) {
    this.serve_date = helper.GetServeDate();

  }

  myproject = false;

  ngOnInit(): void {
    // Check for query parameter
    this.route.queryParams.subscribe(params => {
      this.myproject = params['myproject'] === 'true';
    });
    
    // Google Maps API is automatically loaded by the Angular Google Maps module
    // Just load the project data directly
    this.loadProjectData();

    // Observable to reload data if the user registers
    this.registrationSubscription = this.registrationService.registrationChange$.subscribe(() => {
      this.loadProjectData();
    });
  }

  loadProjectData(): void {
    this.isLoading = true;
    this.isAdmin = this.authService.isAdmin();
    this.userSignedIn = this.authService.isAuthenticated();

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
    this.isLoading = true
    this.projectService.getProject(project_id).subscribe({
      next: (data: Project) => {
        this.project = data;

        // linkify anything in the project description
        const options = {
          target: { url: '_blank' }
        };
        this.project.rich_description = this.linkifyService.linkify(this.project.description, options);

        // if the user is signed in, then we will check to see if they are already registered
        // note - you can hit this page without being signed in
        if (this.userSignedIn) {
          // Check if user is registered for this project
          this.associateProjectToUser(project_id);
        }

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

  associateProjectToUser(project_id: number) {
    this.projectService.getUserRegistrations().subscribe({
      next: (data: Registration[]) => {
        this.isRegistered = data?.some(
            (reg) => reg.project_id === project_id && reg.status === "registered",
        );
      }
    })
  }

  openRegistrationForm(): void {
    this.router.navigate(['/projects', this.project?.id, 'register']);
  }

  registerForProject(data: any): void {
    if (!this.project) return;

    this.loadingRegistration = true;
    this.registrationError = "";

    // Close dialog if open
    if (this.dialogRef) {
      this.dialogRef.close(true);
    }

    let phone_number = data.phone
    if (phone_number.length == 10) {
      phone_number = this.formatPhone(phone_number);
    }

    const body = {
      guest_count: data.guest_count,
      email: data.email,
      phone: phone_number,
      first_name: data.first_name,
      last_name: data.last_name,
      lead_interest: data.lead_interest,
      text_permission: data.text_permission,
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
        if (error.status === 409) {
          this.dialog.open(AlreadyRegisteredElsewhereDialogComponent, {
            width: '400px'
          });
        } else if (error.status === 208) {
          this.dialog.open(AlreadyRegisteredDialogComponent, {
            width: '400px'
          });
        } else {
          this.registrationError =
            error.error?.error || "Failed to register for the project";
          this.helper.showError(this.registrationError);
        }
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

  openEditGuestCountDialog(): void {
    if (!this.project) return;

    this.projectService.getUserRegistrations().subscribe({
      next: (registrations) => {
        const registration = registrations.find(r => r.project_id === this.project?.id);
        if (registration) {
          const dialogRef = this.dialog.open(EditGuestCountDialogComponent, {
            width: '400px',
            data: { registration }
          });

          dialogRef.afterClosed().subscribe(guest_count => {
            if (guest_count !== undefined) {
              // Reload project to get updated capacity
              this.loadProjectDetails(this.project!.id);
            }
          });
        }
      }
    });
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

  formatPhone(phone: string): string {
    // Remove any non-digit characters (optional, for robustness)
    const digits = phone.replace(/\D/g, '');
    // Format as 303-841-6058
    return digits.replace(/(\d{3})(\d{3})(\d{4})/, '$1-$2-$3');
  }
}