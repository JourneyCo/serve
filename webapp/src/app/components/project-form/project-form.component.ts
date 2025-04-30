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
import {ProjectService, GoogleMapsService, UserService, HelperService} from '@services';
import { Project, Tools, Skills, Categories, Ages, Supplies } from "@models";
import { MatSelectModule } from "@angular/material/select";
import {MaterialModule} from '@material';
import {MatCheckbox} from '@angular/material/checkbox';

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
  toolList = Tools;
  toolKeys = Object.keys(Tools);
  skillList = Skills;
  skillKeys = Object.keys(Skills);
  categoryList = Categories;
  categoryKeys = Object.keys(Categories);
  ageList = Ages;
  ageKeys = Object.keys(Ages);
  supplyList = Supplies;
  supplyKeys = Object.keys(Supplies);

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
    this.userService.getAllUsers().subscribe(
        (users) => (this.users = users),
        (error) => console.error("Error loading users:", error),
    );
  }

  getProject(id: number | null) {
    if (id) {
      this.projectService.getProject(id).subscribe(data => {
        this.project = data
        console.log(this.project);
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
      short_description: [
        this.project?.short_description || "",
        [Validators.required, Validators.maxLength(200)],
      ],
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
      location_name: [this.project?.location_name || ""],
      location_address: [this.project?.location_address || ""],
      latitude: [
        this.project?.latitude || null,
        [Validators.pattern(/^-?[0-9]+(\.[0-9]+)?$/)],
      ],
      longitude: [
        this.project?.longitude || null,
        [Validators.pattern(/^-?[0-9]+(\.[0-9]+)?$/)],
      ],
      wheelchair_accessible: [this.project?.wheelchair_accessible || false],
      lead_user_id: [this.project?.lead_user_id || ""],
      tools: [this.project?.tools?.map((t) => t.id) || []],
      supplies: [this.project?.supplies?.map((s) => s.id) || []],
      ages: [this.project?.ages?.map((a) => a.id) || []],
      categories: [this.project?.categories?.map((c) => c.id) || []],
      skills: [this.project?.skills?.map((s) => s.id) || []],
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
      title: formValues.title,
      short_description: formValues.short_description,
      description: formValues.description,
      time: formValues.time,
      max_capacity: formValues.max_capacity,
      current_registrations: this.project?.current_registrations || 0,
      location_name: formValues.location_name || null,
      latitude: formValues.latitude ? Number(formValues.latitude) : null,
      longitude: formValues.longitude ? Number(formValues.longitude) : null,
      wheelchair_accessible: formValues.wheelchair_accessible,
      lead_user_id: formValues.lead_user_id,
      tools: formValues.tools,
      supplies: formValues.supplies,
      ages: formValues.ages,
      categories: formValues.categories,
      skills: formValues.skills,
      location_address: formValues.location_address,
      project_date: formValues.project_date,
      created_at: this.project?.created_at || new Date().toISOString(),
      updated_at: new Date().toISOString(),
    };

    const request = this.data.isEdit
        ? this.projectService.updateProject(project)
        : this.projectService.createProject(project);

    request.subscribe(
        (result) => {
          this.submitting = false;
          this.dialogRef.close(result);
        },
        (error) => {
          this.submitting = false;
          console.error("Error saving project:", error);
          this.helper.showError(error.error?.error || "Failed to save project");
        },
    );
  }
  setupLocationGeocoding(): void {
    this.projectForm
        .get("location_name")
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
        .subscribe(
            (result) => {
              if (result) {
                this.projectForm.patchValue({
                  latitude: result.latitude,
                  longitude: result.longitude,
                });

                this.helper.showSuccess("Location geocoded successfully");
              }
            },
            (error) => {
              console.error("Geocoding error:", error);
              this.helper.showError(
                  "Failed to geocode address: " + (error.message || "Unknown error"),
              );
            },
        );
  }
}
