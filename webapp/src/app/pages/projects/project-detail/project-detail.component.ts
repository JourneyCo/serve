import {Component, Input, OnInit, TemplateRef, ViewChild} from '@angular/core';
import {CommonModule} from '@angular/common';
import {FormsModule} from '@angular/forms';
import {ActivatedRoute, Router, RouterModule} from '@angular/router';
import {MatDialog} from '@angular/material/dialog';
import {GoogleMapsModule} from '@angular/google-maps';
import {AuthService, HelperService, ProjectService, RegistrationService, ServeCookie} from '@services';
import {Observable, Subscription} from 'rxjs';
import {Ages, Project, Registration, User} from '@models';
import {
  EditGuestCountDialogComponent
} from '@components';
import { CancellationSuccessDialogComponent } from '../../../components/dialogs/cancellation-success-dialog/cancellation-success-dialog.component';
import {MaterialModule} from '@material';
import {NgxLinkifyjsModule, NgxLinkifyjsService} from 'ngx-linkifyjs-v2';
import {CookieService} from 'ngx-cookie-service';
import {PhoneNumberPipe} from '../../../services/pipe.service';


@Component({
  selector: "app-project-detail",
  standalone: true,
  imports: [
    CommonModule,
    FormsModule,
    RouterModule,
    GoogleMapsModule,
    MaterialModule,
    NgxLinkifyjsModule,
    PhoneNumberPipe,
  ],
  templateUrl: "./project-detail.component.html",
  styleUrls: ["./project-detail.component.scss"],
})

export class ProjectDetailComponent implements OnInit {
  @Input() proj_id: string = '';
  @Input() admin_route: boolean;

  project: Project | null = null;
  currentUser: User | null = null;
  isRegistered = false;
  isLoading = true;
  loadingRegistration = false;
  registrationError = "";
  types: Record<number, string> = {};
  ages = Ages
  serve_date: Date;
  registrationSubscription: Subscription;
  isAdmin: Observable<boolean>;
  userSignedIn = false;

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
  userEmail: string;
  myproject = false;
  projectID: number

  constructor(
    private route: ActivatedRoute,
    private router: Router,
    private projectService: ProjectService,
    private authService: AuthService,
    private dialog: MatDialog,
    private helper: HelperService,
    private registrationService: RegistrationService,
    private linkifyService: NgxLinkifyjsService,
    private cookieService: CookieService,
    private serveCookie: ServeCookie
  ) {
    this.serve_date = helper.GetServeDate();
    this.myproject = this.router.getCurrentNavigation()?.extras.state?.['myproject'];
    this.userEmail = this.router.getCurrentNavigation()?.extras.state?.['email'];
  }

  ngOnInit(): void {
    this.isAdmin = this.authService.isAdmin();
    this.authService.isAuthenticated().subscribe(data => {
      this.userSignedIn = data;
    });

    // Load types data first
    this.loadTypes();

    // Google Maps API is automatically loaded by the Angular Google Maps module
    // Just load the project data directly
    this.loadProjectData();

    // Observable to reload data if the user registers
    this.registrationSubscription = this.registrationService.registrationChange$.subscribe(() => {
      this.loadProjectData();
    });
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

  loadProjectData(): void {
    this.isLoading = true;

    // Get the project ID from the route
    const project_id = this.route.snapshot.paramMap.get("id") || this.proj_id;
    if (this.cookieService.get("servedayproject") == project_id) {
      this.myproject = true;
    }

    if (!project_id || isNaN(+project_id)) {
      this.helper.showError("Invalid project ID");
      this.router.navigate(["/projects"]);
      return;
    }
    this.projectID = Number(project_id);

    // Load project data
    this.loadProjectDetails(+project_id);
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
        if (this.project.website != null) {
          this.project.website = this.linkifyService.linkify(this.project.website, options);
        }
        this.project.encoded_address = encodeURIComponent(this.project.location_address || "");
        this.project.serve_lead_name = this.project.serve_lead_name? this.project.serve_lead_name : this.project.serve_lead?.first_name + " " +
            this.project.serve_lead?.last_name;
        this.project.serve_lead_email = this.project.serve_lead_email || this.project.serve_lead?.email;

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
    if ((!this.userEmail || this.userEmail == "") && this.currentUser && this.currentUser.email) {
        this.userEmail = this.currentUser?.email
      } else {
      return
    }
    this.projectService.getUserRegistrations(this.userEmail).subscribe({
      next: (data: Registration) => {
        this.isRegistered = data.status === "registered";
      }
    })
  }

  openRegistrationForm(): void {
    this.router.navigate(['/projects', this.project?.id, 'register'], {
      state: { myproject: true, email: this.userEmail}});
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
  console.log(this.userEmail);
    this.projectService.getUserRegistrations(this.userEmail).subscribe({
      next: (registration) => {
        if (registration) {
          const dialogRef = this.dialog.open(EditGuestCountDialogComponent, {
            width: '400px',
            data: {
              registration: registration,
              project: this.project
            }
          });

          dialogRef.afterClosed().subscribe(guest_count => {
            if (guest_count !== undefined) {
              this.projectService.updateUserRegistration(registration.id, { guest_count: guest_count }).subscribe({
                next: (data) => {
                  this.helper.showSuccess("Guest count successfully updated");
                  // Reload project to get updated capacity
                  this.loadProjectDetails(this.project!.id);
                },
                error: (error: any) => {
                  console.error("Error updating guest count:", error);
                  this.helper.showError("Error updating guest count");
                },
              })}
          });
        }
      },
      error: (error: any) => {
        if (error.status === 428) {
          this.helper.showError("Registration not found. Please try again.");
          this.router.navigate(['/']);
        }
      }
    });
  }

  cancelRegistration(): void {
    if (!this.project) return;
    if ((!this.userEmail || this.userEmail == "") && this.currentUser && this.currentUser.email) {
      this.userEmail = this.currentUser?.email
    }
    if (!this.userEmail || this.userEmail == "") {
      this.userEmail = this.cookieService.get("servedayemail")
    }
    if (!this.userEmail || this.userEmail == "") {
      return
    }

    this.loadingRegistration = true;

    this.projectService.cancelRegistration(this.project.id, this.userEmail).subscribe({
      next: () => {
        // @ts-ignore
        this.serveCookie.DeleteProject(this.project.id)
        this.loadingRegistration = false;
        this.isRegistered = false;
        this.myproject = false;
        
        // Show success dialog
        this.dialog.open(CancellationSuccessDialogComponent, {
          width: '400px',
          disableClose: false
        });

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

  isProjectOpen(): boolean {
    return this.project
        ? this.project.status === 'open'
        : false;
  }

  getCapacityPercentage(): number {
    if (!this.project) return 0;
    return (this.project.current_registrations / this.project.max_capacity) * 100;
  }

  getCapacityStatusClass(project: Project): string {
    if (!project) return '';
    const percentage = (project.current_registrations / project.max_capacity) * 100;
    
    if (percentage >= 100) return 'full';
    if (percentage >= 80) return 'almost-full';
    if (percentage >= 50) return 'filling';
    return 'available';
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