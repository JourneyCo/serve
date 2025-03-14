import {Component, inject, OnInit} from '@angular/core';
import {MAT_DIALOG_DATA, MatDialogClose, MatDialogContent, MatDialogTitle} from "@angular/material/dialog";
import {Form, FormBuilder, FormControl, FormGroup, ReactiveFormsModule, Validators} from "@angular/forms";
import {NgIf, CommonModule} from "@angular/common";
import {MatFormFieldModule} from "@angular/material/form-field";

@Component({
  selector: 'register-modal',
  templateUrl: 'register.component.html',
  imports: [MatDialogTitle, MatDialogContent, ReactiveFormsModule, NgIf, MatDialogClose, MatFormFieldModule],
})

export class RegisterDialogComponent implements OnInit {
  data = inject(MAT_DIALOG_DATA);
  private formBuilder = inject(FormBuilder);
  registerForm!: FormGroup;

  ngOnInit(): void {
    this.registerForm = new FormGroup({
      id: new FormControl(this.data.id),
      user_id: new FormControl(this.data.user_id),
      first_name: new FormControl('', [
        Validators.required,
        // Validators.minLength(4),
      ]),
      last_name: new FormControl('', [Validators.required]),
      email: new FormControl('', [Validators.required, Validators.email]),
      qty_enroll: new FormControl('', [Validators.required, Validators.max(this.data.needed)]),
      phoneNumber: new FormControl('', [Validators.required]),
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
