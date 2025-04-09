import { Injectable } from '@angular/core';
import { HttpClient, HttpHeaders } from '@angular/common/http';
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

  registerForProject(project_id: number, body: any): Observable<Registration> {
    const httpOptions = {
      headers: new HttpHeaders({
        'Content-Type': 'application/json'
      })
    };
    return this.http.post<Registration>(`${this.apiUrl}/projects/${project_id}/register`, body, httpOptions);
  }

  cancelRegistration(project_id: number): Observable<any> {
    return this.http.post<any>(`${this.apiUrl}/projects/${project_id}/cancel`, {});
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

  getProjectRegistrations(project_id: number): Observable<Registration[]> {
    return this.http.get<Registration[]>(`${this.apiUrl}/projects/${project_id}/registrations`);
  }
}
