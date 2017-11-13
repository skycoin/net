import { Component, OnInit, ViewEncapsulation, OnDestroy } from '@angular/core';
import { ApiService, Conn, ConnData, ConnsResponse } from '../../service';
import { DataSource } from '@angular/cdk/collections';
import { Observable } from 'rxjs/Observable';
import { Subscription } from 'rxjs/Subscription';
import { Router } from '@angular/router';
import { MatSnackBar } from '@angular/material';
import 'rxjs/add/observable/timer';
import 'rxjs/add/operator/map';

@Component({
  selector: 'app-dashboard',
  templateUrl: './dashboard.component.html',
  styleUrls: ['./dashboard.component.scss'],
  encapsulation: ViewEncapsulation.None
})
export class DashboardComponent implements OnInit, OnDestroy {
  displayedColumns = ['index', 'status', 'key', 'send', 'recv', 'seen'];
  dataSource: ExampleDataSource = null;
  dataSize = 0;
  timer: Subscription = null;

  constructor(private api: ApiService, private snackBar: MatSnackBar, private router: Router) { }
  ngOnInit() {
    this.refresh();
    this.timer = Observable.timer(0, 5000).subscribe(() => {
      this.refresh();
    });
  }
  ngOnDestroy() {
    this.close();
  }
  status(ago: number) {
    const now = new Date().getTime() / 1000;
    return (now - ago) < 180;
  }
  refresh(ev?: Event) {
    if (ev) {
      ev.stopImmediatePropagation();
      ev.stopPropagation();
      ev.preventDefault();
    }
    this.dataSource = new ExampleDataSource(this.api);
    if (ev) {
      this.snackBar.open('Refreshed', 'Dismiss', {
        duration: 3000,
        verticalPosition: 'top',
        extraClasses: ['bg-success']
      });
    }
  }
  openStatus(ev: Event, conn: Conn) {
    if (!conn) {
      this.snackBar.open('Unable to obtain the node state', 'Dismiss', {
        duration: 3000,
        verticalPosition: 'top'
      });
      return;
    }
    this.router.navigate(['node'], { queryParams: { key: conn.key } });
  }
  close() {
    this.timer.unsubscribe();
  }
}
export class ExampleDataSource extends DataSource<any> {
  size = 0;
  constructor(private api: ApiService) {
    super();
  }
  connect(): Observable<Conn[]> {
    return this.api.getAllNode().map((conns: Array<Conn>) => {
      conns.sort((a, b) => {
        if (a.key !== b.key) {
          return a.key.localeCompare(b.key);
        } else {
          if (a.start_time < b.start_time) {
            return 1;
          }
          if (a.start_time > b.start_time) {
            return -1;
          }
          return 0;
        }
      });
      return conns;
    });
  }

  disconnect() { }
}
