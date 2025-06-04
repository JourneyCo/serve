
import { Component } from '@angular/core';
import { CommonModule } from '@angular/common';
import { MatDialogRef } from '@angular/material/dialog';
import { MaterialModule } from '@material';

@Component({
  selector: 'app-cancellation-success-dialog',
  standalone: true,
  imports: [CommonModule, MaterialModule],
  template: `
    <h2 mat-dialog-title>Registration Cancelled</h2>
    <mat-dialog-content>
      <p>Your registration has been successfully cancelled.</p>
    </mat-dialog-content>
    <mat-dialog-actions align="end">
      <button mat-raised-button color="primary" (click)="onOk()">Ok</button>
    </mat-dialog-actions>
  `
})
export class CancellationSuccessDialogComponent {
  constructor(
    private dialogRef: MatDialogRef<CancellationSuccessDialogComponent>
  ) {}

  onOk(): void {
    this.dialogRef.close();
  }
}
