import {Component, inject} from '@angular/core';
import {MAT_DIALOG_DATA, MatDialogContent, MatDialogTitle} from "@angular/material/dialog";
import {FormBuilder, FormControl, FormGroup, ReactiveFormsModule, Validators} from "@angular/forms";

@Component({
  selector: 'register-modal',
  templateUrl: 'register.component.html',
  imports: [MatDialogTitle, MatDialogContent, ReactiveFormsModule],
})


export class RegisterDialogComponent {
  data = inject(MAT_DIALOG_DATA);
  private formBuilder = inject(FormBuilder);

  registerForm = this.formBuilder.group({
    firstName: ['', Validators.required],
    lastName: ['', Validators.required],
    registering: ['', Validators.required],
    phoneNumber: ['', Validators.required],
  });

  onSubmit() {
    console.log("form submitted");
  }

}
