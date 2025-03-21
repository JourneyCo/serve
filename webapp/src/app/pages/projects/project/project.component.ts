import {Component, inject, OnInit} from '@angular/core';
import {MapComponent} from "@components";
import { APIService } from '@services';
import {Location, Project} from "@models";
import {ActivatedRoute, Router} from "@angular/router";
import {Subject} from "rxjs";
import {CommonModule} from "@angular/common";

@Component({
  selector: 'app-project',
  templateUrl: './project.component.html',
  styleUrl: './project.component.css',
  imports: [CommonModule, MapComponent]
})
export class ProjectComponent implements OnInit {

  private APIService = inject(APIService);
  project: Project;
  location: Location;
  projects: Project[] = [];
  locations: Location[] = [];
  eventsSubject: Subject<any> = new Subject<any>();
  private readonly router = inject(ActivatedRoute);

ngOnInit() {
  const id = Number(this.router.snapshot.paramMap.get('id'));
  this.loadProject(id);
  this.loadLocation(id);
}


  loadProject(id: number) {
    this.APIService.getProject(id).subscribe(data => {
      this.project = data;
      this.projects.push(this.project);
    });
  }

  loadLocation(id: number) {
    this.APIService.getLocation(id).subscribe(data => {
      this.location = data;
      this.locations.push(this.location)
    });
  }


}
