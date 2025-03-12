import {AfterViewInit, Component, EventEmitter, inject, Output, signal, ViewChild, WritableSignal} from '@angular/core';
import {MatPaginator, MatPaginatorModule} from '@angular/material/paginator';
import {MatSort, MatSortModule} from '@angular/material/sort';
import {MatTable, MatTableDataSource, MatTableModule} from '@angular/material/table';
import {MatInputModule} from '@angular/material/input';
import {MatFormFieldModule} from '@angular/material/form-field';
import {Location, Project, Registration} from "@models";
import {APIService} from "@services";
import {MapComponent, RegisterDialogComponent} from "@components";
import {CommonModule, DatePipe, NgIf} from "@angular/common";
import {MatDialog} from "@angular/material/dialog";
import {map, Subject, Subscription} from "rxjs";
import {AuthService} from "@auth0/auth0-angular";
import {CodeSnippetComponent} from "../../components/code-snippet.component";

@Component({
  selector: 'projects',
  styleUrl: 'projects.component.css',
  templateUrl: 'projects.component.html',
  imports: [NgIf, MatFormFieldModule, MatInputModule, MatTableModule, MatSortModule, MatPaginatorModule, MapComponent,
  DatePipe, CommonModule, CodeSnippetComponent],
})


export class ProjectsComponent implements AfterViewInit {
  displayedColumns: string[] = ['id', 'name', 'required', 'needed', 'date', 'created_at', 'updated_at', 'register']
  dataSource: MatTableDataSource<Project> = new MatTableDataSource();
  projects: Project[];
  clickedRow: Project | null;
  dialog = inject(MatDialog);
  @ViewChild(MatPaginator) paginator: MatPaginator;
  @ViewChild(MatSort) sort: MatSort;
  @ViewChild(MatTable) table: MatTable<any>;
  locations: Location[];
  locationMap = new Map<number, any>();
  eventsSubject: Subject<any> = new Subject<any>();
  private APIService = inject(APIService);
  private auth = inject(AuthService);
  title = 'Decoded ID Token';
  user$ = this.auth.user$;
  code$ = this.user$.pipe(map((user) => JSON.stringify(user, null, 2)));
  user_id = '';

  constructor() {
    this.loadProjects(true);
    this.loadLocations();
  }

  ngOnInit() {
    this.user$.subscribe(user=> {
      if (user?.sub) {
        this.user_id = user?.sub
      }
    });
  }

  ngAfterViewInit() {
    this.dataSource.paginator = this.paginator;
    this.dataSource.sort = this.sort;
  }

  loadProjects(init: boolean) {
    this.APIService.getProjects().subscribe(data => {
      this.projects = data;
      if (init) {
        this.dataSource = new MatTableDataSource(this.projects);
      } else {
        this.dataSource.data = [...this.projects];
      }
    });
  }

  loadLocations() {
    this.APIService.getLocations().subscribe(data => {
      this.locations = data;
      this.locations.forEach(location => {
        let id = location.id ? location.id : 0;
        this.locationMap.set(id, location);
      })
    });
  }

  applyFilter(event: Event) {
    const filterValue = (event.target as HTMLInputElement).value;
    this.dataSource.filter = filterValue.trim().toLowerCase();

    if (this.dataSource?.paginator) {
      this.dataSource.paginator.firstPage();
    }
  }

  rowClick(row: Project) {
    if (this.clickedRow == row) {
      this.clickedRow = null
      return
    }
    this.clickedRow = row;
    let id = row.location_id ? row.location_id : 0;
    // @ts-ignore
    const clickedLocation = this.locationMap.get(id);
    this.eventsSubject.next(clickedLocation);
  }

  register(evt: any, row: Project) {
    const buttonElement = document.activeElement as HTMLElement; // Get the currently focused element
    buttonElement.blur();

    const data = {
      ...row,
      user_id: this.user_id,
    }

    const dialogRef = this.dialog.open(RegisterDialogComponent, {
      data: data,
      height: '400px',
      width: '700px',
    });

    dialogRef.afterClosed().subscribe(result => {  // result is a form back from dialog
      if (!result) {
        evt.stopPropagation();
        return
      }
      const rawFormValues = result.getRawValue();
      const registration: Registration = {
        id: rawFormValues.id,
        registering: rawFormValues.registering,
        user_id: this.user_id,
      }
      this.APIService.putRegistration(registration).subscribe(data => {
        this.loadProjects(false);
        this.table.renderRows();
        this.paginator.firstPage();
        }
      )
    });
    evt.stopPropagation();
  }
}
