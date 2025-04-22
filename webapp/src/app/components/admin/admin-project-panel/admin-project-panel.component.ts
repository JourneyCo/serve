import {Component, Input, OnInit} from '@angular/core';
import {HelperService, ProjectService} from '@services';
import {Registration} from '@models';
import {MatTable, MatTableDataSource, MatTableModule} from '@angular/material/table';
import {MatPaginator, MatPaginatorModule} from '@angular/material/paginator';
import {MatIcon} from '@angular/material/icon';
import {MatCard, MatCardContent, MatCardHeader, MatCardTitle} from '@angular/material/card';
import {MatDialog} from '@angular/material/dialog';

@Component({
  selector: 'app-admin-project-panel',
  imports: [MatCardContent,MatCardTitle, MatCardHeader, MatTable, MatCard, MatTableModule, MatPaginatorModule, MatIcon],
  templateUrl: './admin-project-panel.component.html',
  styleUrl: './admin-project-panel.component.scss'
})
export class AdminProjectPanelComponent implements OnInit {
  @Input() project_id: number;
  registrationsColumns = ["userName", "projectTitle", "guestCount", "actions"];
  registrations: Registration[] = [];
  registrationsDataSource = new MatTableDataSource<Registration>();
  loadingRegistrations = true;
  processingAction = false;



  constructor(
    private projectService: ProjectService,
    private dialog: MatDialog,
    private helper: HelperService,
  ) {}

  ngOnInit() {
    this.loadRegistrations(this.project_id);
  }

  loadRegistrations(project: number): void {
    this.processingAction = true;
    this.projectService.getProjectRegistrations(project).subscribe({
      next: (registrations) => {
        console.log(registrations);
        this.registrationsDataSource.data = registrations;
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
        this.loadRegistrations(this.project_id);
        this.helper.showSuccess('Registration deleted successfully');
        this.processingAction = false;
      },
      error: (error) => {
        console.error('Error deleting registration:', error);
        this.helper.showError('Failed to delete registration');
        this.processingAction = false;
      }
    });
  }

  editGuestCount(registration: Registration): void {
    const dialogRef = this.dialog.open(MatDialog, {
      width: '300px',
      data: {
        title: 'Edit Guest Count',
        content: 'Enter new guest count:',
        value: registration.guest_count
      }
    });

    dialogRef.afterClosed().subscribe(result => {
      if (result !== undefined && result >= 0) {
        this.processingAction = true;
        this.projectService.updateRegistration(registration.id, { guest_count: result }).subscribe({
          next: () => {
            this.loadRegistrations(this.project_id);
            this.helper.showSuccess('Guest count updated successfully');
            this.processingAction = false;
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
}
