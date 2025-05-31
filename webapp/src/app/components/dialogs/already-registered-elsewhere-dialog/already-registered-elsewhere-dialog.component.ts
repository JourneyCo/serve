
import { Component } from '@angular/core';
import { MatDialogRef } from '@angular/material/dialog';
import { MaterialModule } from '@material';

@Component({
  selector: 'app-already-registered-elsewhere-dialog',
  standalone: true,
  imports: [MaterialModule],
  template: `
    <h2 mat-dialog-title>Already Registered for Another Project</h2>
    <mat-dialog-content>
      <p>You are already registered for a different project. To view your existing registration, please enter your email in the "My Project" link at the top right of the screen.</p>
    </mat-dialog-content>
    <mat-dialog-actions align="end">
      <button mat-button [mat-dialog-close]="true">OK</button>
    </mat-dialog-actions>
  `
})
export class AlreadyRegisteredElsewhereDialogComponent {
  constructor(public dialogRef: MatDialogRef<AlreadyRegisteredElsewhereDialogComponent>) {}
}
