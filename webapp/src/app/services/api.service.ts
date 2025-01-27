import {inject, Injectable} from "@angular/core";
import {HttpClient} from "@angular/common/http";
import {environment} from "../../environments/environment.development";
import {Observable} from "rxjs";


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

}
