import {Component, inject, Input, OnInit} from '@angular/core';
import {MatTab, MatTabGroup} from "@angular/material/tabs";
import {CommonModule} from "@angular/common";
import {Account, Project} from "@models";
import {FormControl, FormGroup, ReactiveFormsModule, Validators} from "@angular/forms";
import {APIService} from "@services";

@Component({
  selector: 'admin',
  imports: [CommonModule, MatTab, MatTabGroup, ReactiveFormsModule],
  templateUrl: './admin.component.html',
  styleUrl: './admin.component.css'
})
export class AdminComponent implements OnInit {
  @Input() leader: Account;
  @Input() project: Project;
  detailsForm!: FormGroup;
  private APIService = inject(APIService);
  accounts: Account[];


  ngOnInit(): void {
    this.detailsForm = new FormGroup({
      name: new FormControl(this.project.name, [Validators.required]),
      short_description: new FormControl(this.project.short_description, [Validators.required]),
      long_description: new FormControl(this.project.long_description, [Validators.required]),
      enabled: new FormControl(this.project.enabled, [Validators.required]),
      start_time: new FormControl(this.project.start_time, [Validators.required]),
      end_time: new FormControl(this.project.end_time, [Validators.required]),
      wheelchair: new FormControl(this.project.wheelchair, [Validators.required]),
      leader_id: new FormControl(this.project.leader_id, [Validators.required]),
      status: new FormControl(this.project.status, [Validators.required])
    });
    this.APIService.getAccounts(this.project.id).subscribe(data => {
      this.accounts = data;
    })
  }

  onSubmit() {
    if (!this.detailsForm.valid) {
      this.detailsForm.markAllAsTouched();
      return
    }
  };

}
