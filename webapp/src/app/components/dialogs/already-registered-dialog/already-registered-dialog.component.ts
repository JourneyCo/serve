
import { Component } from '@angular/core';
import { MatDialogRef } from '@angular/material/dialog';
import { MaterialModule } from '@material';

@Component({
  selector: 'app-already-registered-dialog',
  standalone: true,
  imports: [MaterialModule],
  template: `
    <h2 mat-dialog-title>Already Registered</h2>
    <mat-dialog-content>
      <p>You are already registered for this project. You can view and manage your registration in the "My Projects" section.</p>
    </mat-dialog-content>
    <mat-dialog-actions align="end">
      <button mat-button [mat-dialog-close]="true">OK</button>
    </mat-dialog-actions>
  `
})
export class AlreadyRegisteredDialogComponent {
  constructor(public dialogRef: MatDialogRef<AlreadyRegisteredDialogComponent>) {}
}
