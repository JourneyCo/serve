
import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterModule } from '@angular/router';
import { MatCardModule } from '@angular/material/card';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatIconModule } from '@angular/material/icon';
import { MatTableModule } from '@angular/material/table';
import { MatButtonModule } from '@angular/material/button';
import { MatTooltipModule } from '@angular/material/tooltip';
import { UserService } from '../../services/user.service';
import { ProjectService } from '../../services/project.service';
import { User } from '../../models/user.model';
import { Registration } from '../../models/registration.model';
import { MatTableDataSource } from '@angular/material/table';

@Component({
  selector: 'app-profile',
  standalone: true,
  imports: [
    CommonModule,
    RouterModule,
    MatCardModule,
    MatProgressSpinnerModule,
    MatIconModule,
    MatTableModule,
    MatButtonModule,
    MatTooltipModule
  ],
  templateUrl: './profile.component.html',
  styleUrls: ['./profile.component.scss']
})
export class ProfileComponent implements OnInit {
  loading = true;
  registrationLoading = false;
  error: string | null = null;
  user: User | null = null;
  registrationsDataSource = new MatTableDataSource<Registration>([]);
  registrationsColumns = ['projectTitle', 'startDate', 'endDate', 'location', 'details', 'status', 'actions'];

  constructor(
      private userService: UserService,
      private projectService: ProjectService
  ) {}

  ngOnInit(): void {
    this.loadUserProfile();
    this.loadUserRegistrations();
  }

  private loadUserProfile(): void {
    this.loading = true;
    this.error = null;

    this.userService.getUserProfile().subscribe({
      next: (user) => {
        this.user = user;
        this.loading = false;
      },
      error: (err) => {
        console.error('Error loading profile:', err);
        this.error = 'Failed to load profile. Please try again later.';
        this.loading = false;
      }
    });
  }

  private loadUserRegistrations(): void {
    this.projectService.getUserRegistrations().subscribe({
      next: (registrations) => {
        this.registrationsDataSource.data = registrations;
      },
      error: (err) => {
        console.error('Error loading registrations:', err);
      }
    });
  }

  getDaysUntilStart(startDate: string | undefined): number {
    if (!startDate) return 0;
    const start = new Date(startDate);
    const today = new Date();
    const diffTime = start.getTime() - today.getTime();
    return Math.ceil(diffTime / (1000 * 60 * 60 * 24));
  }

  cancelRegistration(projectId: number): void {
    this.registrationLoading = true;
    this.projectService.cancelRegistration(projectId).subscribe({
      next: () => {
        this.loadUserRegistrations();
        this.registrationLoading = false;
      },
      error: (err) => {
        console.error('Error cancelling registration:', err);
        this.registrationLoading = false;
      }
    });
  }
}
