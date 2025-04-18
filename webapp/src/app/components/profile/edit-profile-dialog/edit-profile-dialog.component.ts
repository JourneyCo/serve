
import { Component, Inject } from '@angular/core';
import { CommonModule } from '@angular/common';
import { MatDialogRef, MAT_DIALOG_DATA, MatDialogModule } from '@angular/material/dialog';
import { MatButtonModule } from '@angular/material/button';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { MatCheckboxModule } from '@angular/material/checkbox';
import { FormBuilder, FormGroup, ReactiveFormsModule } from '@angular/forms';
import { User } from '@models';

@Component({
    selector: 'app-edit-profile-dialog',
    standalone: true,
    imports: [
        CommonModule,
        MatDialogModule,
        MatButtonModule,
        MatFormFieldModule,
        MatInputModule,
        MatCheckboxModule,
        ReactiveFormsModule
    ],
    templateUrl: './edit-profile-dialog.component.html',
    styleUrls: ['./edit-profile-dialog.component.scss']
})
export class EditProfileDialogComponent {
    profileForm: FormGroup;

    constructor(
        private dialogRef: MatDialogRef<EditProfileDialogComponent>,
        @Inject(MAT_DIALOG_DATA) public data: { user: User },
        private fb: FormBuilder
    ) {
        this.profileForm = this.fb.group({
            email: [data.user.email],
            first_name: [data.user.first_name],
            last_name: [data.user.last_name],
            phone: [data.user.phone],
            text_permission: [data.user.text_permission]
        });
    }

    onSubmit(): void {
        if (this.profileForm.valid) {
            this.dialogRef.close(this.profileForm.value);
        }
    }

    onCancel(): void {
        this.dialogRef.close();
    }
}
