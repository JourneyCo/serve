import {Component, Input, OnInit} from '@angular/core';
import { EditGuestCountDialogComponent } from '../dialogs/edit-guest-count-dialog/edit-guest-count-dialog.component';
import {HelperService, ProjectService} from '@services';
import {Project, Registration} from '@models';
import {MatTable, MatTableDataSource} from '@angular/material/table';
import {MatIcon} from '@angular/material/icon';
import {MatCard, MatCardContent, MatCardHeader, MatCardTitle} from '@angular/material/card';
import {MatDialog} from '@angular/material/dialog';
import {RegistrationService} from '@services';
import {Subscription} from 'rxjs';
import {MaterialModule} from '@material';

@Component({
  selector: 'app-admin-project-panel',
  imports: [MatCardContent,MatCardTitle, MatCardHeader, MatTable, MatCard, MatIcon, MaterialModule],
  templateUrl: './admin-project-panel.component.html',
  styleUrl: './admin-project-panel.component.scss'
})
export class AdminProjectPanelComponent implements OnInit {
  @Input() project: Project;
  registrationsColumns = ["userName", "email", "phone", "guestCount", "lead_interest", "actions"];
  registrations: Registration[] = [];
  registrationsDataSource = new MatTableDataSource<Registration>();
  processingAction = false;
  registrationSubscription: Subscription;


  constructor(
    private projectService: ProjectService,
    private dialog: MatDialog,
    private helper: HelperService,
    private registrationChange: RegistrationService
  ) {}

  ngOnInit() {
    console.log(this.project);
    this.loadRegistrations(this.project.id);
    this.registrationSubscription = this.registrationChange.registrationChange$.subscribe(() => {
      this.loadRegistrations(this.project.id);
    });
  }

  loadRegistrations(project: number): void {
    this.processingAction = true;
    this.projectService.getProjectRegistrations(project).subscribe({
      next: (registrations) => {
        this.registrationsDataSource.data = registrations;
        console.log(registrations);
        this.processingAction = false;
      },
      error: (error) => {
        console.error("Error loading registrations:", error);
        this.helper.showError("Failed to load registrations");
        this.processingAction = false;
      }
    });
  }


  deleteRegistration(registration: Registration): void {
    if (!confirm('Are you sure you want to delete this registration? This cannot be undone.')) {
      return;
    }

    this.processingAction = true;
    this.projectService.deleteRegistration(registration.id).subscribe({
      next: () => {
        this.loadRegistrations(this.project.id);
        this.helper.showSuccess('Registration deleted successfully');
        this.processingAction = false;
        this.registrationChange.triggerRegistrationChange();
      },
      error: (error) => {
        console.error('Error deleting registration:', error);
        this.helper.showError('Failed to delete registration');
        this.processingAction = false;
      }
    });
  }

  editGuestCount(registration: Registration): void {
    const dialogRef = this.dialog.open(EditGuestCountDialogComponent, {
      width: '300px',
      data: { registration: registration, project: this.project }
    });

    dialogRef.afterClosed().subscribe(result => {
      if (result !== undefined && result >= 0) {
        this.processingAction = true;
        this.projectService.updateRegistration(registration.id, { guest_count: result }).subscribe({
          next: () => {
            this.loadRegistrations(this.project.id);
            this.loadProject();
            this.helper.showSuccess('Guest count updated successfully');
            this.processingAction = false;
            this.registrationChange.triggerRegistrationChange();
          },
          error: (error) => {
            console.error('Error updating guest count:', error);
            this.helper.showError('Failed to update guest count');
            this.processingAction = false;
          }
        });
      }
    });
  }

  loadProject() {
    this.projectService.getProject(this.project.id).subscribe({
      next: (data) => {
        this.project = data;
      }
    })
  }
}
