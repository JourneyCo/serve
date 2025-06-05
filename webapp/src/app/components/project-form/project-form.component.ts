import { Component, OnInit, Inject } from "@angular/core";
import { CommonModule } from "@angular/common";
import {
  ReactiveFormsModule,
  FormBuilder,
  FormGroup,
  Validators,
  FormsModule,
} from "@angular/forms";
import {
  MatDialogRef,
  MAT_DIALOG_DATA,
} from "@angular/material/dialog";
import { debounceTime, distinctUntilChanged, finalize } from "rxjs/operators";
import { ProjectService, GoogleMapsService, UserService, HelperService } from '@services';
import { Project, Types, Ages } from "@models";
import { MatSelectModule } from "@angular/material/select";
import { MaterialModule } from '@material';

interface DialogData {
  projectID: number | null;
  isEdit: boolean;
}

@Component({
  selector: "app-project-form",
  standalone: true,
  imports: [
    CommonModule,
    ReactiveFormsModule,
    FormsModule,
    MatSelectModule,
    MaterialModule,
  ],
  templateUrl: "./project-form.component.html",
  styleUrls: ["./project-form.component.scss"],
})

export class ProjectFormComponent implements OnInit {
  loading = true;
  projectForm!: FormGroup;
  project: Project | null;
  submitting = false;
  geocoding = false;
  dialogTitle: string;
  minDate: Date;
  users: any[] = [];
  typeList = Types;
  typeKeys = Object.keys(Types);
  ageList = Ages;

  get sortedUsers() {
    return [...this.users].sort((a, b) => a.last_name.localeCompare(b.last_name));
  }


  constructor(
    private fb: FormBuilder,
    private projectService: ProjectService,
    private mapsService: GoogleMapsService,
    private userService: UserService,
    private dialogRef: MatDialogRef<ProjectFormComponent>,
    private helper: HelperService,
    @Inject(MAT_DIALOG_DATA) public data: DialogData,
  ) {
    this.dialogTitle = data.isEdit ? "Edit Project" : "Create New Project";
    this.minDate = new Date();
    this.minDate.setHours(0, 0, 0, 0);
  }

  ngOnInit(): void {
    this.loadUsers();
    this.getProject(this.data.projectID);
    this.setupLocationGeocoding();
  }

  loadUsers(): void {
    this.userService.getAllUsers().subscribe({
      next: (users) => {
        this.users = users;
      },
      error: (error) => {
        console.error("Error loading users:", error);
      }
    });
  }

  getProject(id: number | null) {
    if (id) {
      this.projectService.getProject(id).subscribe(data => {
        this.project = data
        this.loading = false;
        this.initForm();
      })
    } else {
      this.project = null;
      this.loading = false;
      this.initForm();
    }
  }

  initForm(): void {
    const defaultDate = this.helper.GetServeDate();

    this.projectForm = this.fb.group({
      title: [
        this.project?.title || "",
        [Validators.required, Validators.maxLength(100)],
      ],
      website: [this.project?.website || ""],
      description: [
        this.project?.description || "",
        [Validators.required, Validators.minLength(10)],
      ],
      project_date: [defaultDate],
      time: [this.project?.time || "", Validators.required],
      max_capacity: [
        this.project?.max_capacity || 10,
        [Validators.required, Validators.min(1), Validators.max(1000)],
      ],
      area: [this.project?.area || ""],
      location_address: [this.project?.location_address || ""],
      latitude: [
        this.project?.latitude || null,
        [Validators.pattern(/^-?[0-9]+(\.[0-9]+)?$/)],
      ],
      longitude: [
        this.project?.longitude || null,
        [Validators.pattern(/^-?[0-9]+(\.[0-9]+)?$/)],
      ],
      serve_lead_id: [this.project?.serve_lead_id || ""],
      serve_lead_name: [this.project?.serve_lead_name || ""],
      serve_lead_email: [this.project?.serve_lead_email || ""],
      ages: [this.project?.ages || "", Validators.required],
      types: [this.project?.types?.map((c) => c.id) || []],
    });
  }

  onSubmit(): void {
    if (this.projectForm.invalid) {
      return;
    }

    this.submitting = true;

    const formValues = this.projectForm.value;

    const project: Project = {
      id: this.project?.id || 0,
      google_id: this.project?.google_id || 0,
      title: formValues.title,
      website: formValues.website,
      description: formValues.description,
      time: formValues.time,
      max_capacity: formValues.max_capacity,
      current_registrations: this.project?.current_registrations || 0,
      area: formValues.area || null,
      latitude: formValues.latitude ? Number(formValues.latitude) : null,
      longitude: formValues.longitude ? Number(formValues.longitude) : null,
      serve_lead_id: formValues.serve_lead_id,
      serve_lead_name: formValues.serve_lead_name,
      serve_lead_email: formValues.serve_lead_email,
      ages: formValues.ages,
      types: formValues.types,
      location_address: formValues.location_address,
      project_date: formValues.project_date,
      created_at: this.project?.created_at || new Date().toISOString(),
      updated_at: new Date().toISOString(),
    };

    const request = this.data.isEdit
      ? this.projectService.updateProject(project)
      : this.projectService.createProject(project);

    request.subscribe({
      next: (result) => {
        this.submitting = false;
        this.dialogRef.close(result);
      },
      error: (error) => {
        this.submitting = false;
        console.error("Error saving project:", error);
        this.helper.showError(error.error?.error || "Failed to save project");
      },
    });
  }
  setupLocationGeocoding(): void {
    this.projectForm
      .get("area")
      ?.valueChanges.pipe(debounceTime(800), distinctUntilChanged())
      .subscribe((value) => {
        if (value && value.length > 5) {
          this.geocodeLocation(value);
        }
      });
  }

  geocodeLocation(address: string): void {
    if (this.geocoding || !address) {
      return;
    }

    this.geocoding = true;

    this.mapsService
      .geocodeAddressClientSide(address)
      .pipe(
        finalize(() => {
          this.geocoding = false;
        }),
      )
      .subscribe({
        next: (result) => {
          if (result) {
            this.projectForm.patchValue({
              latitude: result.latitude,
              longitude: result.longitude,
            });

            this.helper.showSuccess("Location geocoded successfully");
          }
        },
        error: (error) => {
          console.error("Geocoding error:", error);
          this.helper.showError(
            "Failed to geocode address: " + (error.message || "Unknown error"),
          );
        },
      });
  }
}
