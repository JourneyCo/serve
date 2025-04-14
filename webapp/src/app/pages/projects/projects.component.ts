import {AfterViewInit, Component, inject, OnInit, ViewChild} from '@angular/core';
import {MatPaginator, MatPaginatorModule} from '@angular/material/paginator';
import {MatSort, MatSortModule} from '@angular/material/sort';
import {MatTable, MatTableDataSource, MatTableModule} from '@angular/material/table';
import {MatInputModule} from '@angular/material/input';
import {MatFormFieldModule} from '@angular/material/form-field';
import {Account, Location, Project, Registration} from "@models";
import {APIService} from "@services";
import {MapComponent, RegisterDialogComponent} from "@components";
import {CommonModule, NgIf} from "@angular/common";
import {MatDialog} from "@angular/material/dialog";
import {map, Subject} from "rxjs";
import {Router, RouterLink} from "@angular/router";
import {MatCard, MatCardContent, MatCardHeader, MatCardTitle} from "@angular/material/card";
import {MatIcon} from "@angular/material/icon";
import {AuthService} from "@auth0/auth0-angular";

@Component({
  selector: 'projects',
  styleUrl: 'projects.component.scss',
  templateUrl: 'projects.component.html',
  imports: [NgIf, MatFormFieldModule, MatInputModule, MatTableModule, MatSortModule, MatPaginatorModule, MapComponent,
  CommonModule, MatCard, MatCardContent, MatCardTitle, MatIcon, MatCardHeader, RouterLink],
})


export class ProjectsComponent implements OnInit, AfterViewInit {
  displayedColumns: string[] = ['name', 'required', 'needed', 'register']
  dataSource: MatTableDataSource<Project> = new MatTableDataSource();
  projects: Project[] = [];
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
  private readonly router = inject(Router);
  title = 'Decoded ID Token';
  user$ = this.auth.user$;
  code$ = this.user$.pipe(map((user) => JSON.stringify(user, null, 2)));
  user_id = '';
  registrations: Registration[];
  registrationMap: Map<number, boolean>;
  user: Account;
  isAuthenticated = false;

  constructor(
    private authService: AuthService,
  ) {
  }

  ngOnInit() {
    this.registrationMap = new Map<number, boolean>();
    this.user$.subscribe(user => {
      if (user?.sub) {
        this.user_id = user?.sub
        this.getUser(user?.sub);
      }
      this.loadRegistrations();
    });
    this.auth.isAuthenticated$.subscribe( status => {
      this.isAuthenticated = status;
      }

    )
    this.loadProjects(true);
    this.loadLocations();
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

  getUser(id: string) {
    this.APIService.getAccount(id).subscribe(data => {
      this.user = data;
    })
  }

  loadRegistrations() {
    this.APIService.getRegistrations().subscribe(data => {
      data.filter((registration: Registration) => registration.account_id >= this.user_id).forEach((reg: Registration) => {
        this.registrationMap.set(reg.project_id, true)
      });
    })
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
    const clickedLocation = this.locationMap.get(id);
    this.eventsSubject.next(clickedLocation);
    const rowID = row ? row.id : null;
    this.router.navigate(['/projects/', rowID]);
  }

  register(evt: any, row: Project) {
    const buttonElement = document.activeElement as HTMLElement; // Get the currently focused element
    buttonElement.blur();

    const data = {
      ...row,
      user: this.user,
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
        project_id: rawFormValues.id,
        first: rawFormValues.first_name,
        last: rawFormValues.last_name,
        email: rawFormValues.email,
        cellphone: rawFormValues.phoneNumber,
        qty_enroll: rawFormValues.qty_enroll,
        account_id: this.user_id,
        lead: rawFormValues.lead,
      }
      this.APIService.putRegistration(registration).subscribe(data => {
          this.loadProjects(false);
          this.loadRegistrations();
          this.table.renderRows();
          this.paginator.firstPage();
        }
      )
    });
    evt.stopPropagation();
  }

  login(): void {
    this.auth.loginWithRedirect({
      appState: {
        target: '/projects',
      },
      authorizationParams: {
        prompt: 'login',
      },
    });
  }
}
