
import { Component } from '@angular/core';
import { CommonModule } from '@angular/common';
import { MatDialogRef } from '@angular/material/dialog';
import { MaterialModule } from '@material';
import { Router } from '@angular/router';

@Component({
  selector: 'app-registration-complete-dialog',
  standalone: true,
  imports: [CommonModule, MaterialModule],
  template: `
    <h2 mat-dialog-title>Registration Complete!</h2>
    <mat-dialog-content>
      <p>Thank you for registering! You will receive periodic updates about your project via email and text messages.</p>
    </mat-dialog-content>
    <mat-dialog-actions align="end">
      <button mat-raised-button color="primary" (click)="onOk()">Ok</button>
    </mat-dialog-actions>
  `
})
export class RegistrationCompleteDialogComponent {
  constructor(
    private dialogRef: MatDialogRef<RegistrationCompleteDialogComponent>,
    private router: Router
  ) {}

  onOk(): void {
    this.dialogRef.close();
    this.router.navigate(['/']);
  }
}
