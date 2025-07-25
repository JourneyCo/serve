import { Component, OnInit, Inject } from "@angular/core";
import { CommonModule } from "@angular/common";
import {
  ReactiveFormsModule,
  FormBuilder,
  FormGroup,
  Validators,
  FormsModule,
} from '@angular/forms';
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
import {ProjectService, GoogleMapsService, HelperService} from '@services';
import {Project, Ages, Type, Registration, Lead} from '@models';
import {MatSelectChange, MatSelectModule} from '@angular/material/select';
import {environment} from "../../../../environments/environment";

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
    MatSelectModule
  ],
  templateUrl: "./project-form.component.html",
  styleUrls: ["./project-form.component.scss"],
})
export class ProjectFormComponent implements OnInit {
  projectForm!: FormGroup;
  submitting = false;
  geocoding = false;
  dialogTitle: string;
  agesList = Ages;
  serve_day: string = environment.serveDay
  types: Record<number, string> = {};
  typeEntries: [string, string][];
  allTypes: Type[];
  registrants: Lead[] = [];



  constructor(
      private fb: FormBuilder,
      private projectService: ProjectService,
      private mapsService: GoogleMapsService,
      private dialogRef: MatDialogRef<ProjectFormComponent>,
      private snackBar: MatSnackBar,
      private helper: HelperService,
      @Inject(MAT_DIALOG_DATA) public data: DialogData,
  ) {
    this.dialogTitle = data.isEdit ? "Edit Project" : "Create New Project";
  }

  ngOnInit(): void {
    this.initForm();
    this.setupLocationGeocoding();
    this.loadTypes();
  }


  initForm(): void {
    const project = this.data.project;

    this.projectForm = this.fb.group({
      title: [
        project?.title || "",
        [Validators.required, Validators.maxLength(100)],
      ],
      description: [
        project?.description || "",
        [Validators.required, Validators.minLength(10)],
      ],
      project_date: project?.project_date,
      time: [project?.time || "", Validators.required],
           max_capacity: [
        project?.max_capacity || 10,
        [Validators.required, Validators.min(1), Validators.max(1000)],
      ],
      latitude: [
        project?.latitude || null,
        [Validators.pattern(/^-?[0-9]+(\.[0-9]+)?$/)],
      ],
      longitude: [
        project?.longitude || null,
        [Validators.pattern(/^-?[0-9]+(\.[0-9]+)?$/)],
      ],
      serve_lead_email: [project?.serve_lead_email || "", [Validators.email, Validators.required]],
      serve_lead_name: [project?.serve_lead_name || "", Validators.required],
      location_address: [project?.location_address || "", Validators.required],
      area: [project?.area || "", Validators.required],
      categories: [project?.types?.map(type => type.id) || []],
      ages: [project?.ages || "All Ages", Validators.required],
      leads: [[]],
    });

    if (project) {
      // Load registrants and update form after they're loaded
      this.projectService.getProjectRegistrations(project.id).subscribe({
        next: (data: Registration[]) => {
          this.registrants = data
              .filter(reg => reg.user?.first_name && reg.user?.last_name)
              .map(reg => ({
                name: `${reg.user!.first_name} ${reg.user!.last_name}`,
                phone: reg.user!.phone,
                email: reg.user!.email,
                active: false
              }));

          // Set active status for existing leads
          if (project.leads && project.leads.length > 0) {
            project.leads.forEach((l) => {
              const lead = this.registrants.find(person => person.email === l.email);
              if (lead) {
                lead.active = true;
              }
            });
          }

          // Update form control with initially selected leads
          const selectedLeads = this.registrants.filter(lead => lead.active);
          this.projectForm.patchValue({ leads: selectedLeads });
        },
        error: (error: any) => {
          console.error("Error loading registration details:", error);
          this.helper.showError("Error loading registration details");
        },
      });
    }

  }

  onSubmit(): void {
    if (this.projectForm.invalid) {
      return;
    }

    this.submitting = true;
    const formValues = this.projectForm.value;

    // Filter only active leads for submission
    const selectedLeads = this.registrants.filter(lead => lead.active);

    const project: Project = {
      id: this.data.project?.id || 0,
      title: formValues.title,
      description: formValues.description,
      time: formValues.time,
      max_capacity: formValues.max_capacity,
      current_registrations: this.data.project?.current_registrations || 0,
      latitude: formValues.latitude ? Number(formValues.latitude) : null,
      longitude: formValues.longitude ? Number(formValues.longitude) : null,
      ages: formValues.ages,
      location_address: formValues.location_address,
      project_date: formValues.project_date,
      created_at: this.data.project?.created_at || new Date().toISOString(),
      updated_at: new Date().toISOString(),
      area: formValues.area,
      serve_lead_name: formValues.serve_lead_name,
      serve_lead_email: formValues.serve_lead_email,
      google_id: 0,
      types: formValues.categories,
      leads: selectedLeads,
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
        .get("location_address")
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

  loadTypes(): void {
    this.projectService.getTypes().subscribe({
      next: (typesData) => {
        this.allTypes = typesData
        // Convert array to Record<number, string> format
        this.types = {};
        typesData.forEach(type => {
          this.types[type.id] = type.name;
        });
        this.typeEntries = Object.entries(this.types);
      },
      error: (error: any) => {
        console.error("Error loading types:", error);
        this.helper.showError("Error loading project types");
      }
    });
  }

  // Update the method signature to handle MatSelectChange
  onLeadSelectionChange(event: MatSelectChange): void {
    const selectedLeads = event.value as Lead[];

    // Reset all leads to inactive first
    this.registrants.forEach(lead => lead.active = false);

    // Set active state for selected leads
    selectedLeads.forEach(selected => {
      const lead = this.registrants.find(r => r.email === selected.email);
      if (lead) {
        lead.active = true;
      }
    });

    // Update form control with selected leads
    this.projectForm.patchValue({ leads: selectedLeads });
  }

}
