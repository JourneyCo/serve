import { Component } from '@angular/core';
import { CommonModule } from '@angular/common';
import { MaterialModule } from '@material';
import { RouterModule, Router } from '@angular/router';
import {HelperService, ProjectService} from '@services';
import { MatDialog } from '@angular/material/dialog';
import { FindProjectDialogComponent } from '../dialogs/find-project-dialog/find-project-dialog.component';
import { Observable } from 'rxjs';
import {AuthService} from '@services';
import { User } from '@models';
import {MatToolbar} from '@angular/material/toolbar';

@Component({
  selector: 'app-nav',
  standalone: true,
  imports: [
    CommonModule,
    RouterModule,
    MaterialModule,
    MatToolbar
  ],
  templateUrl: './nav.component.html',
  styleUrls: ['./nav.component.scss']
})

export class NavComponent {
  isAuthenticated$: Observable<boolean>;
  user$: Observable<User | null>;
  isAdmin: Observable<boolean>;
  serve_day: Date = new Date();

  constructor(
    private dialog: MatDialog,
    private router: Router,
    private authService: AuthService,
    private helperService: HelperService,
    private projectService: ProjectService,
  ) {
    this.isAuthenticated$ = this.authService.isAuthenticated();
    this.user$ = this.authService.getCurrentUser();
    this.isAdmin = this.authService.isAdmin();
    this.serve_day = this.helperService.GetServeDate();
  }

  findMyProject(): void {
    const dialogRef = this.dialog.open(FindProjectDialogComponent, {
      width: '400px'
    });

    dialogRef.afterClosed().subscribe(email => {
      if (email) {
        this.projectService.getMyProject(email).subscribe({
          next: (registration) => {
            if (registration && registration.project_id) {
              this.router.navigateByUrl('/', { skipLocationChange: true }).then(() => {
                this.router.navigate(['/projects', registration.project_id], {
                  state: { myproject: true, email: email, reload: Date.now()}
                });
              });
            } else {
              this.helperService.showError('No project found for this email');
            }
          },
          error: (error) => {
            this.helperService.showError('No project found for this email');
          }
        });
      }
    });
  }

  logout(): void {
    this.authService.logout();
  }
}