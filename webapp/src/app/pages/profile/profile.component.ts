import { Component, OnInit } from "@angular/core";
import { CommonModule } from "@angular/common";
import { RouterModule } from "@angular/router";
import { MatProgressSpinnerModule } from "@angular/material/progress-spinner";
import {UserService, ProjectService, AuthService, HelperService} from '@services';
import { User, Registration } from "@models";
import { Observable } from "rxjs";
import {EditProfileDialogComponent} from "./edit-profile-dialog/edit-profile-dialog.component";
import {MatDialog} from "@angular/material/dialog";
import {MaterialModule} from '@material';
import {MatTableDataSource} from '@angular/material/table';

@Component({
  selector: "app-profile",
  standalone: true,
  imports: [
    CommonModule,
    RouterModule,
    MatProgressSpinnerModule,
    MaterialModule
  ],
  templateUrl: "./profile.component.html",
  styleUrls: ["./profile.component.scss"],
})
export class ProfileComponent implements OnInit {
  loading = true;
  registrationLoading = false;
  error: string | null = null;
  user: User | null = null;
  registrationsDataSource = new MatTableDataSource<Registration>([]);
  registrationsColumns = [
    "projectTitle",
    "time",
    "location",
    "details",
    "status",
    "actions",
  ];
  isAdmin: Observable<boolean>;

  constructor(
      private userService: UserService,
      private projectService: ProjectService,
      private authService: AuthService,
      private dialog: MatDialog,
      private helper: HelperService,
  ) {}

  openEditDialog(): void {
    const dialogRef = this.dialog.open(EditProfileDialogComponent, {
      data: { user: this.user }
    });

    dialogRef.afterClosed().subscribe(result => {
      if (result) {
        this.userService.updateProfile(result).subscribe({
          next: (updatedUser) => {
            this.user = updatedUser;
            this.helper.showSuccess('Profile updated successfully');
          },
          error: (error) => {
            console.error('Error updating profile:', error);
            this.helper.showError('Failed to update profile');
          }
        });
      }
    });
  }

  ngOnInit(): void {
    this.loadUserProfile();
    this.isAdmin = this.authService.isAdmin();
  }

  private loadUserProfile(): void {
    this.userService.getUserProfile().subscribe({
      next: (user) => {
        this.user = user;
        this.loadUserRegistrations();
        this.loading = false;
      },
      error: (err) => {
        console.error("Error loading user profile:", err);
        this.error = "Failed to load profile";
        this.loading = false;
      },
    });
    this.loading = true;
    this.error = null;

    this.userService.getUserProfile().subscribe({
      next: (user) => {
        this.user = user;
        this.loading = false;
      },
      error: (err) => {
        console.error("Error loading profile:", err);
        this.error = "Failed to load profile. Please try again later.";
        this.loading = false;
      },
    });
  }

  private loadUserRegistrations(): void {
    this.projectService.getUserRegistrations().subscribe({
      next: (registrations) => {
        this.registrationsDataSource.data = registrations;
      },
      error: (err) => {
        console.error("Error loading registrations:", err);
      },
    });
  }

  getDaysUntilStart(startDate: string | undefined): number {
    if (!startDate) return 0;
    const start = new Date(startDate);
    const today = new Date();
    const diffTime = start.getTime() - today.getTime();
    return Math.ceil(diffTime / (1000 * 60 * 60 * 24));
  }

  cancelRegistration(project_id: number): void {
    this.registrationLoading = true;
    this.projectService.cancelRegistration(project_id).subscribe({
      next: () => {
        this.loadUserRegistrations();
        this.registrationLoading = false;
      },
      error: (err) => {
        console.error("Error cancelling registration:", err);
        this.registrationLoading = false;
      },
    });
  }
}
