
import { Component, Inject } from '@angular/core';
import { MAT_DIALOG_DATA, MatDialogRef } from '@angular/material/dialog';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { MaterialModule } from '@material';
import { Registration } from '@models';

@Component({
  selector: 'app-edit-guest-count-dialog',
  templateUrl: './edit-guest-count-dialog.component.html',
  styleUrls: ['./edit-guest-count-dialog.component.scss'],
  standalone: true,
  imports: [MaterialModule]
})
export class EditGuestCountDialogComponent {
  form: FormGroup;

  constructor(
    private dialogRef: MatDialogRef<EditGuestCountDialogComponent>,
    @Inject(MAT_DIALOG_DATA) public data: { registration: Registration },
    private fb: FormBuilder
  ) {
    this.form = this.fb.group({
      guest_count: [data.registration.guest_count, [Validators.required, Validators.min(0)]]
    });
  }

  onSubmit(): void {
    if (this.form.valid) {
      this.dialogRef.close(this.form.value.guest_count);
    }
  }

  onCancel(): void {
    this.dialogRef.close();
  }
}
