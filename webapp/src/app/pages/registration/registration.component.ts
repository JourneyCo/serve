
import { Component, OnInit } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import { RegistrationCompleteDialogComponent, AlreadyRegisteredDialogComponent,
  AlreadyRegisteredElsewhereDialogComponent } from '@components';
import { CommonModule } from '@angular/common';
import { FormBuilder, FormGroup, ReactiveFormsModule, Validators } from '@angular/forms';
import { ActivatedRoute, Router } from '@angular/router';
import { MaterialModule } from '@material';
import { NgxMaskDirective } from 'ngx-mask';
import { Project, User } from '@models';
import {HelperService, ProjectService, UserService} from '@services';
import {AuthService} from '@auth0/auth0-angular';
import {RecaptchaModule, ReCaptchaV3Service} from 'ng-recaptcha-2';

@Component({
  selector: 'app-registration',
  standalone: true,
  imports: [
    CommonModule,
    ReactiveFormsModule,
    MaterialModule,
    NgxMaskDirective,
      RecaptchaModule,
      ReactiveFormsModule
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
    private auth0: AuthService,
    private dialog: MatDialog,
    private recaptchaV3Service: ReCaptchaV3Service
  ) {
    this.registrationForm = this.fb.group({
      email: ['', Validators.email],
      first_name: [''],
      last_name: [''],
      guest_count: [0],
      phone: ['', Validators.pattern(/^\d{10}\s*$/)],
      text_permission: [false],
      lead_interest: [false],
      // recaptcha: [null, Validators.required]
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
      // this.recaptchaV3Service.execute('submitRegistration').subscribe({
      //     next: (token) => {
      //       console.log(this.registrationForm.value);
            if (this.registrationForm.valid && this.project) {
              this.projectService.registerForProject(this.project.id, this.registrationForm.value).subscribe({
                next: (response) => {
                  if (response.status === 208) {
                    this.dialog.open(AlreadyRegisteredDialogComponent, {
                      width: '400px',
                      disableClose: true
                    });
                  } else if (response.status === 201 || response.status === 200) {
                    this.dialog.open(RegistrationCompleteDialogComponent, {
                      width: '400px',
                      disableClose: true
                    });
                  }
                },
                error: (error: any) => {
                  if (error.status === 409) {
                    this.dialog.open(AlreadyRegisteredElsewhereDialogComponent, {
                      width: '400px',
                      disableClose: true
                    });
                  } else {
                    console.error('Error registering:', error);
                    this.helper.showError('Error registering for project');
                  }
                },
              });
            }
          // },
          // error: (error: any) => {
          //     console.error('Error procesing captcha:', error);
          //     this.helper.showError('Error registering for project');
          // },
      // });
  }

  onCancel(): void {
    this.router.navigate(['/projects', this.project?.id]);
  }
}
