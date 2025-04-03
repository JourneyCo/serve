import { Component, OnInit, Inject } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ReactiveFormsModule, FormBuilder, FormGroup, Validators } from '@angular/forms';
import { MatDialogRef, MAT_DIALOG_DATA, MatDialogModule } from '@angular/material/dialog';
import { MatSnackBar, MatSnackBarModule } from '@angular/material/snack-bar';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { MatDatepickerModule } from '@angular/material/datepicker';
import { MatNativeDateModule } from '@angular/material/core';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatTooltipModule } from '@angular/material/tooltip';
import { debounceTime, distinctUntilChanged, finalize } from 'rxjs/operators';
import { ProjectService } from '../../../services/project.service';
import { GoogleMapsApiService } from '../../../services/google-maps-api.service';
import { Project } from '../../../models/project.model';

interface DialogData {
  project: Project | null;
  isEdit: boolean;
}

@Component({
  selector: 'app-project-form',
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
    MatTooltipModule
  ],
  templateUrl: './project-form.component.html',
  styleUrls: ['./project-form.component.scss']
})
export class ProjectFormComponent implements OnInit {
  projectForm!: FormGroup;
  submitting = false;
  geocoding = false;
  dialogTitle: string;
  minDate: Date;

  constructor(
    private fb: FormBuilder,
    private projectService: ProjectService,
    private mapsService: GoogleMapsApiService,
    private dialogRef: MatDialogRef<ProjectFormComponent>,
    private snackBar: MatSnackBar,
    @Inject(MAT_DIALOG_DATA) public data: DialogData
  ) {
    this.dialogTitle = data.isEdit ? 'Edit Project' : 'Create New Project';
    
    // Set minimum date for start date (today)
    this.minDate = new Date();
    this.minDate.setHours(0, 0, 0, 0);
  }

  ngOnInit(): void {
    this.initForm();
    this.setupLocationGeocoding();
  }

  initForm(): void {
    const project = this.data.project;
    
    // Default values for time
    const defaultDate = new Date('2025-07-12');
    const defaultStartTime = '09:00';
    const defaultEndTime = '17:00';
    
    this.projectForm = this.fb.group({
      title: [project?.title || '', [Validators.required, Validators.maxLength(100)]],
      description: [project?.description || '', [Validators.required, Validators.minLength(10)]],
      projectDate: [project?.projectDate ? new Date(project.projectDate) : defaultDate, Validators.required],
      startTime: [project?.startTime || defaultStartTime, Validators.required],
      endTime: [project?.endTime || defaultEndTime, Validators.required],
      maxCapacity: [project?.maxCapacity || 10, [Validators.required, Validators.min(1), Validators.max(1000)]],
      locationName: [project?.locationName || ''],
      latitude: [project?.latitude || null, [Validators.pattern(/^-?[0-9]+(\.[0-9]+)?$/)]],
      longitude: [project?.longitude || null, [Validators.pattern(/^-?[0-9]+(\.[0-9]+)?$/)]]
    }, { validators: this.timeRangeValidator });
  }

  timeRangeValidator(formGroup: FormGroup): { [key: string]: boolean } | null {
    const startTime = formGroup.get('startTime')?.value;
    const endTime = formGroup.get('endTime')?.value;
    
    if (startTime && endTime && endTime <= startTime) {
      return { 'timeRange': true };
    }
    
    return null;
  }

  onSubmit(): void {
    if (this.projectForm.invalid) {
      return;
    }
    
    this.submitting = true;
    
    // Format the form values
    const formValues = this.projectForm.value;
    
    // Format the project date to YYYY-MM-DD
    const projectDate = this.formatDate(formValues.projectDate);
    
    // Create project object
    const project: Project = {
      id: this.data.project?.id || 0,
      title: formValues.title,
      description: formValues.description,
      projectDate,
      startTime: formValues.startTime,
      endTime: formValues.endTime,
      maxCapacity: formValues.maxCapacity,
      currentRegistrations: this.data.project?.currentRegistrations || 0,
      locationName: formValues.locationName || null,
      latitude: formValues.latitude ? Number(formValues.latitude) : null,
      longitude: formValues.longitude ? Number(formValues.longitude) : null,
      createdAt: this.data.project?.createdAt || new Date().toISOString(),
      updatedAt: new Date().toISOString()
    };
    
    // Submit the project
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
        console.error('Error saving project:', error);
        this.showError(error.error?.error || 'Failed to save project');
      }
    );
  }

  formatDate(date: Date): string {
    const d = new Date(date);
    let month = '' + (d.getMonth() + 1);
    let day = '' + d.getDate();
    const year = d.getFullYear();

    if (month.length < 2) month = '0' + month;
    if (day.length < 2) day = '0' + day;

    return [year, month, day].join('-');
  }

  private showError(message: string): void {
    this.snackBar.open(message, 'Close', {
      duration: 5000,
      panelClass: ['error-snackbar']
    });
  }

  /**
   * Set up location geocoding to automatically convert addresses to lat/lng
   */
  setupLocationGeocoding(): void {
    // Watch for changes to the locationName field
    this.projectForm.get('locationName')?.valueChanges
      .pipe(
        debounceTime(800), // Wait for the user to stop typing
        distinctUntilChanged() // Only proceed if the value actually changed
      )
      .subscribe(value => {
        // Only attempt geocoding if there's a value with enough characters
        if (value && value.length > 5) {
          this.geocodeLocation(value);
        }
      });
  }

  /**
   * Geocode a location string to latitude/longitude
   */
  geocodeLocation(address: string): void {
    // Don't geocode if already geocoding or if address is empty
    if (this.geocoding || !address) {
      return;
    }

    this.geocoding = true;
    
    this.mapsService.geocodeAddress(address)
      .pipe(
        finalize(() => {
          this.geocoding = false;
        })
      )
      .subscribe(
        result => {
          if (result) {
            // Update the form with the geocoded coordinates
            this.projectForm.patchValue({
              latitude: result.latitude,
              longitude: result.longitude
            });
            
            this.snackBar.open('Location geocoded successfully', 'Close', {
              duration: 3000
            });
          }
        },
        error => {
          console.error('Geocoding error:', error);
          this.showError('Failed to geocode address: ' + (error.message || 'Unknown error'));
        }
      );
  }
}
