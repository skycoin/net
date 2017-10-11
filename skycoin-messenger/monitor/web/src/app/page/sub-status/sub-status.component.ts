import { Component, OnInit, ViewEncapsulation, OnDestroy } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { ApiService, NodeServices } from '../../service';
import { MatSnackBar } from '@angular/material';

@Component({
  selector: 'app-sub-status',
  templateUrl: './sub-status.component.html',
  styleUrls: ['./sub-status.component.scss'],
  encapsulation: ViewEncapsulation.None
})
export class SubStatusComponent implements OnInit, OnDestroy {
  key = '';
  status: NodeServices = null;
  task = null;
  constructor(
    private router: Router,
    private route: ActivatedRoute,
    private api: ApiService,
    private snackBar: MatSnackBar) { }

  ngOnInit() {
    this.route.queryParams.subscribe(params => {
      this.key = params.key;
      this.init();
      this.task = setInterval(() => {
        this.init();
      }, 10000);
    });
  }
  ngOnDestroy() {
    this.close();
  }
  refresh(ev: Event) {
    ev.stopImmediatePropagation();
    ev.stopPropagation();
    ev.preventDefault();
    this.init();
    this.snackBar.open('Refresh Successful', 'Dismiss', {
      duration: 3000,
      verticalPosition: 'top'
    });
  }
  close() {
    clearInterval(this.task);
  }
  init() {
    const data = new FormData();
    data.append('key', this.key);
    this.api.getNodeStatus(data).subscribe((nodeServices: NodeServices) => {
      if (nodeServices) {
        this.status = nodeServices;
      }
    });
  }
}
