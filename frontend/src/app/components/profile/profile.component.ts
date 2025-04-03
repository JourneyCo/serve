import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterModule } from '@angular/router';
import { MatCardModule } from '@angular/material/card';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatTableModule, MatTableDataSource } from '@angular/material/table';
import { MatSnackBar, MatSnackBarModule } from '@angular/material/snack-bar';
import { MatTooltipModule } from '@angular/material/tooltip';
import { AuthService } from '../../services/auth.service';
import { ProjectService } from '../../services/project.service';
import { User } from '../../models/user.model';
import { Registration } from '../../models/registration.model';

@Component({
  selector: 'app-profile',
  standalone: true,
  imports: [
    CommonModule,
    RouterModule,
    MatCardModule,
    MatButtonModule,
    MatIconModule,
    MatProgressSpinnerModule,
    MatTableModule,
    MatSnackBarModule,
    MatTooltipModule
  ],
  templateUrl: './profile.component.html',
  styleUrls: ['./profile.component.scss']
})
export class ProfileComponent implements OnInit {
  user: User | null = null;
  isAdmin = false;
  isLoading = true;
  registrationLoading = false;
  
  // Registrations table
  registrationsDataSource = new MatTableDataSource<Registration>([]);
  registrationsColumns: string[] = ['projectTitle', 'startDate', 'endDate', 'location', 'details', 'status', 'actions'];
  
  constructor(
    private authService: AuthService,
    private projectService: ProjectService,
    private snackBar: MatSnackBar
  ) {}

  ngOnInit(): void {
    this.loadUserData();
  }

  loadUserData(): void {
    this.isLoading = true;
    
    this.authService.getCurrentUser().subscribe(
      user => {
        this.user = user;
        this.isAdmin = user?.isAdmin || false;
        this.loadRegistrations();
      },
      error => {
        console.error('Error loading user profile:', error);
        this.showError('Failed to load user profile');
        this.isLoading = false;
      }
    );
  }

  loadRegistrations(): void {
    this.projectService.getUserRegistrations().subscribe(
      registrations => {
        this.registrationsDataSource.data = registrations;
        this.isLoading = false;
      },
      error => {
        console.error('Error loading registrations:', error);
        this.showError('Failed to load your registrations');
        this.isLoading = false;
      }
    );
  }

  cancelRegistration(projectId: number): void {
    if (!confirm('Are you sure you want to cancel this registration?')) {
      return;
    }
    
    this.registrationLoading = true;
    
    this.projectService.cancelRegistration(projectId).subscribe(
      () => {
        this.registrationLoading = false;
        this.showSuccess('Registration cancelled successfully');
        this.loadRegistrations(); // Reload registrations
      },
      error => {
        this.registrationLoading = false;
        console.error('Error cancelling registration:', error);
        this.showError(error.error?.error || 'Failed to cancel registration');
      }
    );
  }

  formatDate(date: string): string {
    return new Date(date).toLocaleDateString('en-US', { 
      year: 'numeric', 
      month: 'long', 
      day: 'numeric' 
    });
  }

  getDaysUntilStart(startDate: string): number {
    const today = new Date();
    today.setHours(0, 0, 0, 0);
    
    const start = new Date(startDate);
    start.setHours(0, 0, 0, 0);
    
    const diffTime = start.getTime() - today.getTime();
    return Math.ceil(diffTime / (1000 * 60 * 60 * 24));
  }

  private showSuccess(message: string): void {
    this.snackBar.open(message, 'Close', {
      duration: 3000,
      panelClass: ['success-snackbar']
    });
  }

  private showError(message: string): void {
    this.snackBar.open(message, 'Close', {
      duration: 5000,
      panelClass: ['error-snackbar']
    });
  }
}
