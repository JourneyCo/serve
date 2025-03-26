import {Component, inject, OnInit} from '@angular/core';
import {MapComponent} from "@components";
import { APIService } from '@services';
import {Location, Project, Account } from "@models";
import {ActivatedRoute} from "@angular/router";
import {Subject} from "rxjs";
import {CommonModule} from "@angular/common";
import {AuthService} from "@auth0/auth0-angular";
import { jwtDecode } from 'jwt-decode';
import {AdminComponent} from "@components";

@Component({
  selector: 'app-project',
  templateUrl: './project.component.html',
  styleUrl: './project.component.css',
  imports: [CommonModule, MapComponent, AdminComponent]
})
export class ProjectComponent implements OnInit {

  private APIService = inject(APIService);
  project: Project;
  location: Location;
  projects: Project[] = [];
  locations: Location[] = [];
  eventsSubject: Subject<any> = new Subject<any>();
  private readonly router = inject(ActivatedRoute);
  directionURL: string;
  leader: Account;
  leadermailto: string;
  private auth = inject(AuthService);
  user$ = this.auth.user$;
  canAdmin: boolean = false;

ngOnInit() {
  const id = Number(this.router.snapshot.paramMap.get('id'));
  this.loadProject(id);
  this.loadLocation(id);
  this.getUser()
}

  loadProject(id: number) {
    this.APIService.getProject(id).subscribe(data => {
      this.project = data;
      this.projects.push(this.project);
      this.getLeader(this.project.leader_id)
      console.log(this.project);
    });
  }

  loadLocation(id: number) {
    this.APIService.getLocation(id).subscribe(data => {
      this.location = data;
      this.locations.push(this.location);
      const encoded = encodeURIComponent(this.location.formatted_address);
      this.directionURL = "https://www.google.com/maps/dir/?api=1&destination=" + encoded
    });
  }

  getLeader(id: string) {
    this.APIService.getAccount(id).subscribe(data => {
      this.leader = data;
      this.leadermailto = "mailto:" + this.leader.email;
    })
  }

  getUser() {
    this.auth.getAccessTokenSilently().subscribe(token => {
      this.getPermissionsFromToken(token);
    });
  }

  getPermissionsFromToken(accessToken: string) {
      const decodedToken: any = jwtDecode(accessToken);
      const perms =  decodedToken.permissions || [];
      if (perms.includes('edit:projects')) {
        this.canAdmin = true;
      }
    }
}
