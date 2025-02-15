import {inject, Injectable} from "@angular/core";
import {HttpClient, HttpHeaders} from "@angular/common/http";
import {environment} from "../../environments/environment.development";
import {Observable} from "rxjs";
import {Registration} from "@models";


@Injectable()
export class APIService {
  private apiUrl = environment.apiUrl;
  private readonly http: HttpClient = inject(HttpClient)


    public getLocations(): Observable<any[]> {
      return this.http.get<any[]>(`${this.apiUrl}/locations`);
    }

  public getProjects(): Observable<any[]> {
    return this.http.get<any[]>(`${this.apiUrl}/projects`);
  }

  public getSystem(): Observable<any> {
    return this.http.get<any>(`${this.apiUrl}/system`);
  }

  public putRegistration(r: Registration): Observable<any> {
    let headers = new HttpHeaders();
    headers = headers.set('Content-Type', 'application/json');
    const body = JSON.stringify(r);
    const s = r.id.toString();
    return this.http.put<any>(`${this.apiUrl}/projects/`+ s + '/register', body, { headers: headers});
  }

}
