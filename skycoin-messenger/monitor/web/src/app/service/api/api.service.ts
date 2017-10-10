import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { environment } from '../../../environments/environment';
import { Observable } from 'rxjs/Observable';
import 'rxjs/add/operator/catch';
import 'rxjs/add/observable/throw';

@Injectable()
export class ApiService {
  private url = 'http://127.0.0.1:5998/';
  private connUrl = this.url + 'conn/';
  constructor(private httpClient: HttpClient) {
    if (environment.production) {
      this.connUrl = '/conn/';
    }
  }
  getAllNode() {
    return this.handleGet(this.connUrl + 'getAll');
  }

  getNodeStatus(data: FormData) {
    return this.handlePost(this.connUrl + 'getNodeStatus', data);
  }

  handleGet(url: string) {
    if (url === '') {
      return Observable.throw('Url is empty.');
    }
    return this.httpClient.get(url).catch(err => Observable.throw(err));
  }
  handlePost(url: string, data: FormData) {
    if (url === '') {
      return Observable.throw('Url is empty.');
    }
    return this.httpClient.post(url, data).catch(err => Observable.throw(err));
  }
}
export interface Conn {
  key?: string;
  type?: string;
  send_bytes?: number;
  recv_bytes?: number;
  last_ack_time?: number;
  start_time?: number;
}
export interface ConnsResponse {
  conns?: Array<Conn>;
}

export interface NodeServices extends Conn {
  apps?: Array<string>;
}

export interface Services {
  key?: string;
  attributes?: Array<string>;
  address?: string;
}
