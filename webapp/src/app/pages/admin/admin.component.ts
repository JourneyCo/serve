import {AfterViewInit, Component, OnInit, ViewChild} from '@angular/core';
import {CommonModule} from '@angular/common';
import {RouterLink} from '@angular/router';
import {MatDialog, MatDialogModule} from '@angular/material/dialog';
import {MatTabGroup, MatTabsModule} from '@angular/material/tabs';
import {MatTableDataSource, MatTableModule} from '@angular/material/table';
import {MatSort, MatSortModule} from '@angular/material/sort';
import {MatSnackBar, MatSnackBarModule} from '@angular/material/snack-bar';
import {MatButtonModule} from '@angular/material/button';
import {MatIconModule} from '@angular/material/icon';
import {MatCardModule} from '@angular/material/card';
import {MatInputModule} from '@angular/material/input';
import {MatFormFieldModule} from '@angular/material/form-field';
import {MatProgressSpinnerModule} from '@angular/material/progress-spinner';
import {MatDividerModule} from '@angular/material/divider';
import {MatTooltipModule} from '@angular/material/tooltip';
import {MatSlideToggleModule} from '@angular/material/slide-toggle';
import {ProjectService} from '@services';
import {Project} from '@models';
import {ProjectFormComponent} from '../../components/admin/project-form/project-form.component';

@Component({
  selector: "app-admin-dashboard",
  standalone: true,
  imports: [
    CommonModule,
    RouterLink,
    MatDialogModule,
    MatTabsModule,
    MatTableModule,
    MatSortModule,
    MatSnackBarModule,
    MatButtonModule,
    MatIconModule,
    MatCardModule,
    MatInputModule,
    MatFormFieldModule,
    MatProgressSpinnerModule,
    MatDividerModule,
    MatTooltipModule,
    MatSlideToggleModule,
  ],
  templateUrl: "./admin.component.html",
  styleUrls: ["./admin.component.scss"],
})
export class AdminComponent implements OnInit, AfterViewInit {
  // Data sources for tables
  projectsDataSource = new MatTableDataSource<Project>([]);

  // Column definitions
  projectColumns: string[] = [
    "title",
    "time",
    "capacity",
    "active",
    "actions",
  ];

  // Loading states
  loadingProjects = true;
  processingAction = false;
  loadingStats = true;

  // Statistics
  totalRegistrants = 0;
  totalGuests = 0;

  // ViewChild references for table sorting
  @ViewChild("projectsSort") projectsSort!: MatSort;
  @ViewChild("usersSort") usersSort!: MatSort;
  @ViewChild("tabGroup") tabGroup!: MatTabGroup;

  constructor(
    private projectService: ProjectService,
    private dialog: MatDialog,
    private snackBar: MatSnackBar,
  ) {}

  ngOnInit(): void {
    this.loadProjects();
    this.loadRegistrationStats();
  }

  ngAfterViewInit(): void {
    // Set up sorting after view init
    setTimeout(() => {
      this.projectsDataSource.sort = this.projectsSort;
    });
  }

  loadProjects(): void {
    this.loadingProjects = true;
    this.projectService.getProjects().subscribe(
      (data) => {
        this.projectsDataSource.data = this.sortProjects(data);
        this.loadingProjects = false;
      },
      (error) => {
        console.error("Error loading projects:", error);
        this.showError("Failed to load projects");
        this.loadingProjects = false;
      },
    );
  }

  applyProjectFilter(event: Event): void {
    const filterValue = (event.target as HTMLInputElement).value;
    this.projectsDataSource.filter = filterValue.trim().toLowerCase();
  }

  createProject(): void {
    const dialogRef = this.dialog.open(ProjectFormComponent, {
      width: "600px",
      data: { project: null, isEdit: false },
    });

    dialogRef.afterClosed().subscribe((result) => {
      if (result) {
        this.loadProjects();
        this.showSuccess("Project created successfully");
      }
    });
  }

  editProject(project: Project): void {
    const dialogRef = this.dialog.open(ProjectFormComponent, {
      width: "600px",
      data: { project: { ...project }, isEdit: true },
    });

    dialogRef.afterClosed().subscribe((result) => {
      if (result) {
        this.loadProjects();
        this.showSuccess("Project updated successfully");
      }
    });
  }

  deleteProject(project: Project): void {
    if (
      !confirm(
        `Are you sure you want to delete the project: "${project.title}"? This cannot be undone.`,
      )
    ) {
      return;
    }

    this.processingAction = true;
    this.projectService.deleteProject(project.id).subscribe(
      () => {
        this.loadProjects();
        this.processingAction = false;
        this.showSuccess("Project deleted successfully");
      },
      (error) => {
        console.error("Error deleting project:", error);
        this.showError("Failed to delete project");
        this.processingAction = false;
      },
    );
  }

  toggleProjectActiveStatus(project: Project, event: any): void {
    const newStatus = event.checked ? 'active' : 'inactive';
    
    this.processingAction = true;
    this.projectService.updateProjectActiveStatus(project.id, newStatus).subscribe(
      () => {
        project.active = event.checked;
        this.processingAction = false;
        this.showSuccess(`Project ${newStatus === 'active' ? 'activated' : 'deactivated'} successfully`);
      },
      (error) => {
        console.error("Error updating project status:", error);
        this.showError("Failed to update project status");
        this.processingAction = false;
        // Revert the toggle state on error
        event.source.checked = !event.checked;
      },
    );
  }

  formatDate(date: string): string {
    return new Date(date).toLocaleDateString();
  }

  isProjectAtCapacity(project: Project): boolean {
    return project.current_registrations >= project.max_capacity;
  }

  private showSuccess(message: string): void {
    this.snackBar.open(message, "Close", {
      duration: 3000,
      panelClass: ["success-snackbar"],
    });
  }

  private showError(message: string): void {
    this.snackBar.open(message, "Close", {
      duration: 5000,
      panelClass: ["error-snackbar"],
    });
  }

  loadRegistrationStats(): void {
    this.loadingStats = true;
    this.projectService.getAllRegistrations().subscribe(
      (registrations) => {
        this.totalRegistrants = registrations.length;
        this.totalGuests = registrations.reduce((sum, reg) => sum + (reg.guest_count || 0), 0);
        this.loadingStats = false;
      },
      (error) => {
        console.error("Error loading registration stats:", error);
        this.loadingStats = false;
      }
    );
  }

  sortProjects(projects: Project[]) {
    return projects.sort((a, b) => {
      const aFull = a.current_registrations >= a.max_capacity;
      const bFull = b.current_registrations >= b.max_capacity;

      // If both projects are either full or not full, sort alphabetically by name
      if (aFull === bFull) {
        return a.title.localeCompare(b.title);
      }

      // Projects that are full go to the bottom
      return aFull ? 1 : -1;
    });
  }
}
