import {Component, Inject} from '@angular/core';
import {FormBuilder, FormGroup, ReactiveFormsModule, Validators} from '@angular/forms';
import {MAT_DIALOG_DATA, MatDialogRef} from '@angular/material/dialog';
import {Project, User} from '@models';
import {MaterialModule} from '@material';
import {CommonModule} from '@angular/common';
import {NgxMaskDirective} from 'ngx-mask';

@Component({
  selector: 'app-registration-dialog',
  imports: [MaterialModule, ReactiveFormsModule, CommonModule, NgxMaskDirective],
  templateUrl: './registration-dialog.component.html',
  styleUrl: './registration-dialog.component.scss'
})
export class RegistrationDialogComponent {
  registrationForm: FormGroup;

  constructor(
    private dialogRef: MatDialogRef<RegistrationDialogComponent>,
    @Inject(MAT_DIALOG_DATA) public data: { project: Project, user: User },
    private fb: FormBuilder
  ) {
    this.registrationForm = this.fb.group({
      email: [data.user.email || '', Validators.email],
      first_name: [data.user.first_name || ''],
      last_name: [data.user.last_name || ''],
      guest_count: [0, Validators.max(this.data.project.max_capacity-data.project.current_registrations - 1)],
      phone: [data.user.phone || '', Validators.pattern(/^\d{3}-\d{3}-\d{4}\s*$/)],
      text_permission: [data.user.text_permission || false],
      lead_interest: [data.user.text_permission || false]
    });
  }

  onSubmit(): void {
    if (this.registrationForm.valid) {
      const result = { success: true, values: this.registrationForm.value}
      this.dialogRef.close(result);
    }
  }

  onCancel(): void {
    this.dialogRef.close();
  }

}
