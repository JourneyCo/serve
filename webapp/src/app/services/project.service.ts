import { Injectable } from "@angular/core";
import {HttpClient, HttpHeaders, HttpResponse} from '@angular/common/http';
import { Observable } from "rxjs";
import { environment } from "../../environments/environment";
import { Project, Registration } from '@models';

@Injectable({
  providedIn: "root",
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

  registerForProject(project_id: number, body: any): Observable<HttpResponse<Registration>> {
    const httpOptions = {
      headers: new HttpHeaders({
        "Content-Type": "application/json",
      }),
      observe: 'response' as const
    };
    return this.http.post<Registration>(
      `${this.apiUrl}/projects/${project_id}/register`,
      body,
      httpOptions,
    );
  }

  cancelRegistration(project_id: number, email: string): Observable<any> {
    return this.http.post<any>(
      `${this.apiUrl}/projects/${project_id}/cancel`, {},
        { params: { email } },
    );
  }

  getUserRegistrations(email: string): Observable<Registration> {
    return this.http.get<Registration>(`${this.apiUrl}/users/registrations`, { params: { email } });
  }

  // Admin API endpoints
  createProject(project: Project): Observable<Project> {
    return this.http.post<Project>(`${this.apiUrl}/admin/projects`, project);
  }

  updateProject(project: Project): Observable<Project> {
    return this.http.put<Project>(
      `${this.apiUrl}/admin/projects/${project.id}`,
      project,
    );
  }

  deleteProject(id: number): Observable<any> {
    return this.http.delete<any>(`${this.apiUrl}/admin/projects/${id}`);
  }

  getProjectRegistrations(project_id: number): Observable<Registration[]> {
    return this.http.get<Registration[]>(
      `${this.apiUrl}/projects/${project_id}/registrations`,
    );
  }

  getMyProject(email: string): Observable<Registration> {
    return this.http.get<Registration>(
      `${this.apiUrl}/projects/my`,
      { params: { email } }
    );
  }

  getTypes(): Observable<any[]> {
    return this.http.get<any[]>(`${this.apiUrl}/projects/types`);
  }

  updateRegistration(id: number, updates: { guest_count: number }): Observable<any> {
    return this.http.put(`${this.apiUrl}/admin/registrations/${id}`, updates);
  }

  updateUserRegistration(id: number, updates: { guest_count: number }): Observable<any> {
    return this.http.put(`${this.apiUrl}/users/registrations/${id}`, updates);
  }

  deleteRegistration(id: number): Observable<any> {
    return this.http.delete(`${this.apiUrl}/admin/registrations/${id}`);
  }

  getAllRegistrations(): Observable<Registration[]> {
    return this.http.get<Registration[]>(`${this.apiUrl}/admin/registrations`);
  }
}
