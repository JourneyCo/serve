
import { Component } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormBuilder, FormGroup, ReactiveFormsModule, Validators } from '@angular/forms';
import { MatDialogRef } from '@angular/material/dialog';
import { MaterialModule } from '@material';

@Component({
  selector: 'app-find-project-dialog',
  standalone: true,
  imports: [CommonModule, MaterialModule, ReactiveFormsModule],
  template: `
    <h2 mat-dialog-title>Find My Project</h2>
    <form [formGroup]="findProjectForm" (ngSubmit)="onSubmit()">
      <mat-dialog-content>
        <mat-form-field appearance="outline" style="width: 100%">
          <mat-label>Email Address</mat-label>
          <input matInput formControlName="email" type="email" placeholder="Enter your email">
          <mat-error *ngIf="findProjectForm.get('email')?.hasError('required')">Email is required</mat-error>
          <mat-error *ngIf="findProjectForm.get('email')?.hasError('email')">Please enter a valid email</mat-error>
        </mat-form-field>
      </mat-dialog-content>
      <mat-dialog-actions align="end">
        <button mat-button mat-dialog-close>Cancel</button>
        <button mat-raised-button color="primary" type="submit" [disabled]="findProjectForm.invalid">Find Project</button>
      </mat-dialog-actions>
    </form>
  `
})
export class FindProjectDialogComponent {
  findProjectForm: FormGroup;

  constructor(
    private fb: FormBuilder,
    private dialogRef: MatDialogRef<FindProjectDialogComponent>
  ) {
    this.findProjectForm = this.fb.group({
      email: ['', [Validators.required, Validators.email]]
    });
  }

  onSubmit(): void {
    if (this.findProjectForm.valid) {
      this.dialogRef.close(this.findProjectForm.value.email);
    }
  }
}
