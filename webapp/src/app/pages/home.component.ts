import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { Router, RouterModule } from '@angular/router';
import { MatDialog } from '@angular/material/dialog';
import { AuthService, HelperService } from '@services';
import { MaterialModule } from '@material';
import { FindProjectDialogComponent } from '../components/dialogs/find-project-dialog/find-project-dialog.component';

@Component({
  selector: 'app-home',
  standalone: true,
  imports: [
    CommonModule,
    RouterModule,
    MaterialModule,
  ],
  templateUrl: './home.component.html',
  styleUrls: ['./home.component.scss']
})
export class HomeComponent implements OnInit {
  isAuthenticated = false;
  serve_day: Date = new Date();

  constructor(
    private authService: AuthService,
    private router: Router,
    private helperService: HelperService,
    private dialog: MatDialog,
    private projectService: ProjectService
  ) {
      this.serve_day = this.helperService.GetServeDate();
  }

  ngOnInit(): void {
    this.authService.isAuthenticated().subscribe(
      (isAuthenticated) => {
        this.isAuthenticated = isAuthenticated;
      }
    );
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
              this.router.navigate(['/projects', registration.project_id]);
            }
          },
          error: (error) => {
            this.helperService.showError('No project found for this email');
          }
        });
      }
    });
  }
}
