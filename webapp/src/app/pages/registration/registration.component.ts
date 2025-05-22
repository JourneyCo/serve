
import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormBuilder, FormGroup, ReactiveFormsModule, Validators } from '@angular/forms';
import { ActivatedRoute, Router } from '@angular/router';
import { MaterialModule } from '@material';
import { NgxMaskDirective } from 'ngx-mask';
import { Project, User } from '@models';
import {HelperService, ProjectService, UserService} from '@services';
import {AuthService} from '@auth0/auth0-angular';

@Component({
  selector: 'app-registration',
  standalone: true,
  imports: [
    CommonModule,
    ReactiveFormsModule,
    MaterialModule,
    NgxMaskDirective
  ],
  templateUrl: './registration.component.html',
  styleUrls: ['./registration.component.scss']
})
export class RegistrationComponent implements OnInit {
  registrationForm: FormGroup;
  project: Project;
  user: User | null = null;

  constructor(
    private fb: FormBuilder,
    private route: ActivatedRoute,
    private router: Router,
    private projectService: ProjectService,
    private userService: UserService,
    private helper: HelperService,
    private auth0: AuthService
  ) {
    this.registrationForm = this.fb.group({
      email: ['', Validators.email],
      first_name: [''],
      last_name: [''],
      guest_count: [0],
      phone: ['', Validators.pattern(/^\d{10}\s*$/)],
      text_permission: [false],
      lead_interest: [false]
    });
  }

  ngOnInit() {
    const projectId = this.route.snapshot.paramMap.get('id');

    if (projectId) {
      this.projectService.getProject(Number(projectId)).subscribe(project => {
        this.project = project;
        if (this.auth0.isAuthenticated$) {
          this.userService.getUserProfile().subscribe(user => {
            this.user = user;
            this.registrationForm.patchValue({
              email: user.email || '',
              first_name: user.first_name || '',
              last_name: user.last_name || '',
              phone: user.phone || '',
              text_permission: user.text_permission || false
            });
          });
        }
      });
    }
  }

  onSubmit(): void {
    if (this.registrationForm.valid && this.project) {
      this.projectService.registerForProject(this.project.id, this.registrationForm.value)
        .subscribe(() => {
          this.router.navigate(['/projects', this.project?.id]);
        });

      this.projectService.registerForProject(this.project.id, this.registrationForm.value).subscribe({
        next: (user) => {
          this.router.navigate(['/projects', this.project?.id]);
          this.helper.showSuccess("You have sucessfully registered!");
        },
        error: (error: any) => {
          console.error("Error registering:", error);
          this.helper.showError("Error registering for project");
        },
      });


    }
  }

  onCancel(): void {
    this.router.navigate(['/projects', this.project?.id]);
  }
}
