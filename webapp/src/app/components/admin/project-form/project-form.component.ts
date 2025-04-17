import { Component, OnInit, Inject } from "@angular/core";
import { CommonModule } from "@angular/common";
import {
  ReactiveFormsModule,
  FormBuilder,
  FormGroup,
  Validators,
  FormArray,
  NgModel,
  FormsModule,
  FormControl,
} from "@angular/forms";
import {
  MatDialogRef,
  MAT_DIALOG_DATA,
  MatDialogModule,
} from "@angular/material/dialog";
import { MatSnackBar, MatSnackBarModule } from "@angular/material/snack-bar";
import { MatButtonModule } from "@angular/material/button";
import { MatIconModule } from "@angular/material/icon";
import { MatFormFieldModule } from "@angular/material/form-field";
import { MatInputModule } from "@angular/material/input";
import { MatDatepickerModule } from "@angular/material/datepicker";
import { MatNativeDateModule, MatOption } from "@angular/material/core";
import { MatProgressSpinnerModule } from "@angular/material/progress-spinner";
import { MatTooltipModule } from "@angular/material/tooltip";
import { MatCheckboxModule } from "@angular/material/checkbox";
import { MatChipsModule } from "@angular/material/chips";
import { debounceTime, distinctUntilChanged, finalize } from "rxjs/operators";
import { ProjectService } from "../../../services/project.service";
import { GoogleMapsApiService } from "../../../services/google-maps-api.service";
import { UserService } from "../../../services/user.service";
import { Project } from "../../../models/project.model";
import {
  Tools,
  Skills,
  ProjectAccessory,
  Categories,
  Ages,
  Supplies,
} from "../../../models/accessories";
import { MatSelectModule } from "@angular/material/select";

interface DialogData {
  project: Project | null;
  isEdit: boolean;
}

@Component({
  selector: "app-project-form",
  standalone: true,
  imports: [
    CommonModule,
    ReactiveFormsModule,
    MatDialogModule,
    MatSnackBarModule,
    MatButtonModule,
    MatIconModule,
    MatFormFieldModule,
    MatInputModule,
    MatDatepickerModule,
    MatNativeDateModule,
    MatProgressSpinnerModule,
    MatTooltipModule,
    MatCheckboxModule,
    MatChipsModule,
    MatOption,
    FormsModule,
    MatSelectModule,
  ],
  templateUrl: "./project-form.component.html",
  styleUrls: ["./project-form.component.scss"],
})
export class ProjectFormComponent implements OnInit {
  projectForm!: FormGroup;
  submitting = false;
  geocoding = false;
  dialogTitle: string;
  minDate: Date;
  users: any[] = [];
  get sortedUsers() {
    return [...this.users].sort((a, b) => a.last_name.localeCompare(b.last_name));
  }
  toolList = Tools;
  toolKeys = Object.keys(Tools);
  tool_list: any;
  skillList = Skills;
  skillKeys = Object.keys(Skills);
  skill_list: any;
  categoryList = Categories;
  categoryKeys = Object.keys(Categories);
  category_list: any;
  ageList = Ages;
  ageKeys = Object.keys(Ages);
  age_list: any;
  supplyList = Supplies;
  supplyKeys = Object.keys(Supplies);
  supply_list: any;

  constructor(
      private fb: FormBuilder,
      private projectService: ProjectService,
      private mapsService: GoogleMapsApiService,
      private userService: UserService,
      private dialogRef: MatDialogRef<ProjectFormComponent>,
      private snackBar: MatSnackBar,
      @Inject(MAT_DIALOG_DATA) public data: DialogData,
  ) {
    this.dialogTitle = data.isEdit ? "Edit Project" : "Create New Project";
    this.minDate = new Date();
    this.minDate.setHours(0, 0, 0, 0);
  }

  ngOnInit(): void {
    this.loadUsers();
    this.initForm();
    this.setupLocationGeocoding();
  }

  loadUsers(): void {
    this.userService.getAllUsers().subscribe(
        (users) => (this.users = users),
        (error) => console.error("Error loading users:", error),
    );
  }

  initForm(): void {
    const project = this.data.project;
    console.log(project);

    const defaultDate = "2025-07-12";

    this.projectForm = this.fb.group({
      title: [
        project?.title || "",
        [Validators.required, Validators.maxLength(100)],
      ],
      short_description: [
        project?.short_description || "",
        [Validators.required, Validators.maxLength(200)],
      ],
      description: [
        project?.description || "",
        [Validators.required, Validators.minLength(10)],
      ],
      project_date: [defaultDate, Validators.required],
      time: [project?.time || "", Validators.required],
      max_capacity: [
        project?.max_capacity || 10,
        [Validators.required, Validators.min(1), Validators.max(1000)],
      ],
      location_name: [project?.location_name || ""],
      latitude: [
        project?.latitude || null,
        [Validators.pattern(/^-?[0-9]+(\.[0-9]+)?$/)],
      ],
      longitude: [
        project?.longitude || null,
        [Validators.pattern(/^-?[0-9]+(\.[0-9]+)?$/)],
      ],
      wheelchair_accessible: [project?.wheelchair_accessible || false],
      lead_user_id: [project?.lead_user_id || ""],
      tools: [project?.tools?.map((t) => t.id) || []],
      supplies: [project?.supplies?.map((s) => s.id) || []],
      ages: [project?.ages?.map((a) => a.id) || []],
      categories: [project?.categories?.map((c) => c.id) || []],
      skills: [project?.skills?.map((s) => s.id) || []],
    });
  }

  onSubmit(): void {
    if (this.projectForm.invalid) {
      return;
    }

    this.submitting = true;

    const formValues = this.projectForm.value;

    const project: Project = {
      id: this.data.project?.id || 0,
      title: formValues.title,
      short_description: formValues.short_description,
      description: formValues.description,
      time: formValues.time,
      max_capacity: formValues.max_capacity,
      current_registrations: this.data.project?.current_registrations || 0,
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
      created_at: this.data.project?.created_at || new Date().toISOString(),
      updated_at: new Date().toISOString(),
    };

    console.log(project);
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
          this.showError(error.error?.error || "Failed to save project");
        },
    );
  }

  formatDate(date: Date): string {
    const d = new Date(date);
    let month = "" + (d.getMonth() + 1);
    let day = "" + d.getDate();
    const year = d.getFullYear();

    if (month.length < 2) month = "0" + month;
    if (day.length < 2) day = "0" + day;

    return [year, month, day].join("-");
  }

  private showError(message: string): void {
    this.snackBar.open(message, "Close", {
      duration: 5000,
      panelClass: ["error-snackbar"],
    });
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
        .geocodeAddress(address)
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

                this.snackBar.open("Location geocoded successfully", "Close", {
                  duration: 3000,
                });
              }
            },
            (error) => {
              console.error("Geocoding error:", error);
              this.showError(
                  "Failed to geocode address: " + (error.message || "Unknown error"),
              );
            },
        );
  }
}
