import {AfterViewInit, Component, OnInit, ViewChild} from '@angular/core';
import { CommonModule } from "@angular/common";
import { RouterLink } from "@angular/router";
import { MatDialog } from "@angular/material/dialog";
import { MatTabGroup } from "@angular/material/tabs";
import { MatTableDataSource } from "@angular/material/table";
import { MatPaginator } from "@angular/material/paginator";
import { MatSort } from "@angular/material/sort";
import {HelperService, ProjectService} from '@services';
import {Project, Registration} from '@models';
import { ProjectFormComponent } from '@components';
import {MaterialModule} from '@material';

@Component({
  selector: "app-admin",
  standalone: true,
  imports: [
    CommonModule,
    RouterLink,
    MaterialModule,
  ],
  templateUrl: "./admin.component.html",
  styleUrls: ["./admin.component.scss"],
})
export class AdminComponent implements OnInit, AfterViewInit {
  // Data sources for tables
  projectsDataSource = new MatTableDataSource<Project>([]);

  // Column definitions
  projectColumns: string[] = [
    "id", "google_id",
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
      data: { projectID: null, isEdit: false },
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
      data: { projectID: project.id, isEdit: true },
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
