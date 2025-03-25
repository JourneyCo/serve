import {Component, inject, OnInit} from '@angular/core';
import {MAT_DIALOG_DATA, MatDialogClose, MatDialogContent, MatDialogTitle} from "@angular/material/dialog";
import {FormControl, FormGroup, ReactiveFormsModule, Validators} from "@angular/forms";
import {CommonModule, NgIf} from "@angular/common";
import {MatFormFieldModule} from "@angular/material/form-field";

@Component({
  selector: 'register-modal',
  templateUrl: 'register.component.html',
  imports: [CommonModule, MatDialogTitle, MatDialogContent, ReactiveFormsModule, NgIf, MatDialogClose, MatFormFieldModule],
})

export class RegisterDialogComponent implements OnInit {
  data = inject(MAT_DIALOG_DATA);
  registerForm!: FormGroup;

  ngOnInit(): void {
    this.registerForm = new FormGroup({
      id: new FormControl(this.data.id),
      user_id: new FormControl(this.data.user_id),
      first_name: new FormControl(this.data.user.first, [Validators.required]),
      last_name: new FormControl(this.data.user.last, [Validators.required]),
      email: new FormControl(this.data.user.email, [Validators.required, Validators.email]),
      qty_enroll: new FormControl('', [Validators.required, Validators.max(this.data.required - this.data.registered)]),
      phoneNumber: new FormControl(this.data.user.cellphone, [Validators.required]),
      lead: new FormControl(false),
    });
  }

  onSubmit() {
    if (!this.registerForm.valid) {
      this.registerForm.markAllAsTouched();
      return
    }
  };
}
