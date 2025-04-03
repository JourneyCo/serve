import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { environment } from '../../environments/environment';
import { Project } from '../models/project.model';
import { Registration } from '../models/registration.model';

@Injectable({
  providedIn: 'root'
})
export class ProjectService {
  private apiUrl = `${environment.apiUrl}/api`;

  constructor(private http: HttpClient) {}

  // Public API endpoints
  getProjects(): Observable<Project[]> {
    return this.http.get<Project[]>(`${this.apiUrl}/projects`);
  }

  getProject(id: number): Observable<Project> {
    return this.http.get<Project>(`${this.apiUrl}/projects/${id}`);
  }

  registerForProject(projectId: number, guestCount: number = 0, isProjectLead: boolean = false): Observable<Registration> {
    return this.http.post<Registration>(`${this.apiUrl}/projects/${projectId}/register`, {
      guestCount,
      isProjectLead
    });
  }

  cancelRegistration(projectId: number): Observable<any> {
    return this.http.post<any>(`${this.apiUrl}/projects/${projectId}/cancel`, {});
  }

  getUserRegistrations(): Observable<Registration[]> {
    return this.http.get<Registration[]>(`${this.apiUrl}/users/registrations`);
  }

  // Admin API endpoints
  createProject(project: Project): Observable<Project> {
    return this.http.post<Project>(`${this.apiUrl}/admin/projects`, project);
  }

  updateProject(project: Project): Observable<Project> {
    return this.http.put<Project>(`${this.apiUrl}/admin/projects/${project.id}`, project);
  }

  deleteProject(id: number): Observable<any> {
    return this.http.delete<any>(`${this.apiUrl}/admin/projects/${id}`);
  }

  getProjectRegistrations(projectId: number): Observable<Registration[]> {
    return this.http.get<Registration[]>(`${this.apiUrl}/projects/${projectId}/registrations`);
  }
}
