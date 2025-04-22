import {AfterViewInit, Component, OnInit, ViewChild} from '@angular/core';
import { CommonModule } from "@angular/common";
import { RouterLink } from "@angular/router";
import { MatDialog, MatDialogModule } from "@angular/material/dialog";
import { MatTabsModule, MatTabGroup } from "@angular/material/tabs";
import { MatTableModule, MatTableDataSource } from "@angular/material/table";
import { MatPaginatorModule, MatPaginator } from "@angular/material/paginator";
import { MatSortModule, MatSort } from "@angular/material/sort";
import { MatSnackBarModule } from "@angular/material/snack-bar";
import { MatButtonModule } from "@angular/material/button";
import { MatIconModule } from "@angular/material/icon";
import { MatCardModule } from "@angular/material/card";
import { MatInputModule } from "@angular/material/input";
import { MatFormFieldModule } from "@angular/material/form-field";
import { MatProgressSpinnerModule } from "@angular/material/progress-spinner";
import { MatDividerModule } from "@angular/material/divider";
import { MatTooltipModule } from "@angular/material/tooltip";
import {HelperService, ProjectService} from '@services';
import {Project, Registration} from '@models';
import { ProjectFormComponent } from "../project-form/project-form.component";

@Component({
  selector: "app-admin-dashboard",
  standalone: true,
  imports: [
    CommonModule,
    RouterLink,
    MatDialogModule,
    MatTabsModule,
    MatTableModule,
    MatPaginatorModule,
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
  ],
  templateUrl: "./admin-dashboard.component.html",
  styleUrls: ["./admin-dashboard.component.scss"],
})
export class AdminDashboardComponent implements OnInit, AfterViewInit {
  // Data sources for tables
  projectsDataSource = new MatTableDataSource<Project>([]);

  // Column definitions
  projectColumns: string[] = [
    "title",
    "time",
    "capacity",
    "actions",
  ];

  registrationsDataSource = new MatTableDataSource<Registration>([]);

  // Loading states
  loadingProjects = true;
  processingAction = false;

  // ViewChild references for table sorting and pagination
  @ViewChild("projectsPaginator") projectsPaginator!: MatPaginator;
  @ViewChild("projectsSort") projectsSort!: MatSort;
  @ViewChild("usersPaginator") usersPaginator!: MatPaginator;
  @ViewChild("usersSort") usersSort!: MatSort;
  @ViewChild("tabGroup") tabGroup!: MatTabGroup;

  constructor(
    private projectService: ProjectService,
    private dialog: MatDialog,
    private helper: HelperService

  ) {}

  ngOnInit(): void {
    this.loadProjects();
    this.loadRegistrations();
  }

  @ViewChild('registrationsPaginator') registrationsPaginator!: MatPaginator;

  ngAfterViewInit(): void {
    // Set up sorting and pagination after view init
    setTimeout(() => {
      this.projectsDataSource.paginator = this.projectsPaginator;
      this.projectsDataSource.sort = this.projectsSort;
      this.registrationsDataSource.paginator = this.registrationsPaginator;
    });
  }

  loadProjects(): void {
    this.loadingProjects = true;
    this.projectService.getProjects().subscribe(
      (projects) => {
        this.projectsDataSource.data = projects;
        this.loadingProjects = false;
      },
      (error) => {
        console.error("Error loading projects:", error);
        this.helper.showError("Failed to load projects");
        this.loadingProjects = false;
      },
    );
  }

  loadRegistrations(): void {
    this.processingAction = true;
    this.projectService.getAllRegistrations().subscribe({
      next: (registrations) => {
        this.registrationsDataSource.data = registrations;
        this.processingAction = false;
      },
      error: (error) => {
        console.error("Error loading registrations:", error);
        this.helper.showError("Failed to load registrations");
        this.processingAction = false;
      }
    });
  }

  applyProjectFilter(event: Event): void {
    const filterValue = (event.target as HTMLInputElement).value;
    this.projectsDataSource.filter = filterValue.trim().toLowerCase();

    if (this.projectsDataSource.paginator) {
      this.projectsDataSource.paginator.firstPage();
    }
  }

  createProject(): void {
    const dialogRef = this.dialog.open(ProjectFormComponent, {
      width: "600px",
      data: { project: null, isEdit: false },
    });

    dialogRef.afterClosed().subscribe((result) => {
      if (result) {
        this.loadProjects();
        this.helper.showSuccess("Project created successfully");
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
        this.helper.showSuccess("Project updated successfully");
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
        this.helper.showSuccess("Project deleted successfully");
      },
      (error) => {
        console.error("Error deleting project:", error);
        this.helper.showError("Failed to delete project");
        this.processingAction = false;
      },
    );
  }

  isProjectAtCapacity(project: Project): boolean {
    return project.current_registrations >= project.max_capacity;
  }
}
