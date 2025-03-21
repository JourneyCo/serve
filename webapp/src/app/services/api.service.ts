import {Injectable} from "@angular/core";
import {HttpClient, HttpHeaders} from "@angular/common/http";
import {environment} from "../../environments/environment.development";
import {Observable} from "rxjs";
import {Registration} from "@models";


@Injectable({
  providedIn: 'root'
})
export class APIService {
  private apiUrl = environment.apiUrl;

  constructor(private http: HttpClient) {
  }

  public getLocation(id: number): Observable<any> {
    return this.http.get<any[]>(`${this.apiUrl}/locations/`+id);
  }

  public getLocations(): Observable<any[]> {
      return this.http.get<any[]>(`${this.apiUrl}/locations`);
    }

  public getProjects(): Observable<any[]> {
    return this.http.get<any[]>(`${this.apiUrl}/projects`);
  }

  public getProject(id: number): Observable<any> {
    return this.http.get<any[]>(`${this.apiUrl}/projects/`+id);
  }

  public getRegistrations(): Observable<any> {
    return this.http.get<any>(`${this.apiUrl}/registrations`);
  }

  public putRegistration(r: Registration): Observable<any> {
    let headers = new HttpHeaders();
    headers = headers.set('Content-Type', 'application/json');
    const body = JSON.stringify(r);
    return this.http.put<any>(`${this.apiUrl}/projects/`+ r.project_id + '/register', body, { headers: headers});
  }

}
