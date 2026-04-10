// src/app/app.component.ts
import { Component, ChangeDetectorRef } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { HttpClient } from '@angular/common/http';

interface Alert {
  type: 'success' | 'error' | 'warning' | 'info';
  message: string;
  id: number;
}

interface UploadedFile {
  name: string;
  size: string;
  uploadTime: string;
  status: 'success' | 'pending';
}

@Component({
  selector: 'app-root',
  standalone: true,
  imports: [CommonModule, FormsModule],
  template: `
    <div [ngSwitch]="isLoggedIn">
      <div *ngSwitchCase="false">
        <ng-container *ngTemplateOutlet="loginTemplate"></ng-container>
      </div>
      <div *ngSwitchDefault>
        <ng-container *ngTemplateOutlet="mainTemplate"></ng-container>
      </div>
    </div>

    <!-- LOGIN TEMPLATE -->
    <ng-template #loginTemplate>
      <div class="container" style="margin-top: 5rem;">
        <div class="card" style="max-width: 400px; margin: 0 auto;">
          <div style="text-align: center; margin-bottom: 2rem;">
            <i class="fas fa-cloud-upload" style="font-size: 3rem; color: var(--primary-color); margin-bottom: 1rem; display: block;"></i>
            <h1 style="color: var(--text-primary); margin-bottom: 0.5rem;">File Upload Manager</h1>
            <p style="color: var(--text-secondary);">Secure file management system</p>
          </div>

          <div class="form-group">
            <label class="form-label">Email</label>
            <input 
              type="email" 
              class="form-input" 
              [(ngModel)]="credentials.email"
              placeholder="you@example.com"
            />
          </div>

          <div class="form-group">
            <label class="form-label">Password</label>
            <input 
              type="password" 
              class="form-input" 
              [(ngModel)]="credentials.password"
              placeholder="Enter your password"
              (keyup.enter)="login()"
            />
          </div>

          <button (click)="login()" class="btn btn-primary btn-block" [disabled]="isLoading">
            <span *ngIf="isLoading" class="spinner"></span>
            <span *ngIf="!isLoading"><i class="fas fa-sign-in-alt"></i></span>
            {{ isLoading ? 'Logging in...' : 'Login' }}
          </button>

          <p style="text-align: center; margin-top: 1rem; color: var(--text-secondary); font-size: 0.875rem;">
            <strong>Enter your credentials to login</strong><br/>
            email<br/>
            password<br/>
          </p>
        </div>
      </div>
    </ng-template>

    <!-- MAIN TEMPLATE -->
    <ng-template #mainTemplate>
      <div class="container">
        <!-- Navigation -->
        <div class="navbar">
          <div class="navbar-brand">
            <i class="fas fa-cloud-upload" style="margin-right: 0.5rem;"></i>FileVault
          </div>
          <nav class="navbar-nav">
            <a class="navbar-link" [class.active]="activeTab === 'upload'" (click)="activeTab = 'upload'">
              <i class="fas fa-upload"></i> Upload
            </a>
            <a class="navbar-link" [class.active]="activeTab === 'history'" (click)="activeTab = 'history'">
              <i class="fas fa-history"></i> History
            </a>
          </nav>
          <div class="user-info">
            <div class="user-avatar">{{ userInitial }}</div>
            <span style="color: var(--text-primary); font-weight: 500;">{{ userEmail }}</span>
            <button class="logout-btn" (click)="logout()">
              <i class="fas fa-sign-out-alt"></i> Logout
            </button>
          </div>
        </div>

        <!-- Alerts -->
        <div *ngFor="let alert of alerts" [ngClass]="'alert alert-' + alert.type">
          <i class="fas" 
            [ngClass]="{
              'fa-check-circle alert-icon': alert.type === 'success',
              'fa-exclamation-circle alert-icon': alert.type === 'error',
              'fa-info-circle alert-icon': alert.type === 'info'
            }"></i>
          <span>{{ alert.message }}</span>
          <button class="alert-close" (click)="removeAlert(alert.id)">
            <i class="fas fa-times"></i>
          </button>
        </div>

        <!-- Main Content -->
        <div class="main-content">
          <!-- Upload Section -->
          <div *ngIf="activeTab === 'upload'" class="card">
            <h2 style="margin-bottom: 2rem; color: var(--text-primary);">
              <i class="fas fa-file-upload"></i> Upload Files
            </h2>

            <!-- Drag & Drop Area -->
            <div 
              class="upload-area"
              [class.drag-over]="isDragging"
              [class.has-file]="selectedFile"
              (click)="fileInput.click()"
              (dragover)="onDragOver($event)"
              (dragleave)="onDragLeave($event)"
              (drop)="onDrop($event)"
            >
              <input 
                #fileInput 
                type="file" 
                class="file-input-hidden" 
                (change)="onFileSelected($event)"
              />
              <div class="upload-icon">
                <i class="fas" [ngClass]="selectedFile ? 'fa-file-check' : 'fa-cloud-upload-alt'"></i>
              </div>
              <div class="upload-text" *ngIf="!selectedFile">
                Drag and drop your file here
              </div>
              <div class="upload-text" *ngIf="selectedFile">
                {{ selectedFile.name }}
              </div>
              <div class="upload-subtext">
                {{ selectedFile ? 'Click to change file' : 'or click to browse' }}
              </div>
            </div>

            <!-- File Info -->
            <div *ngIf="selectedFile" class="file-info">
              <div class="file-icon">
                <i class="fas fa-file"></i>
              </div>
              <div class="file-details">
                <div class="file-name">{{ selectedFile.name }}</div>
                <div class="file-size">{{ (selectedFile.size / 1024).toFixed(2) }} KB</div>
              </div>
              <button class="remove-file" (click)="clearFile()">
                <i class="fas fa-trash"></i> Clear
              </button>
            </div>

            <!-- Progress Bar -->
            <div *ngIf="uploadProgress > 0 && uploadProgress < 100" class="progress-container">
              <div class="progress-label">
                <span>Uploading...</span>
                <span>{{ uploadProgress }}%</span>
              </div>
              <div class="progress-bar">
                <div class="progress-fill" [style.width.%]="uploadProgress"></div>
              </div>
            </div>

            <!-- Upload Button -->
            <button 
              (click)="upload()" 
              class="btn btn-primary btn-block"
              [disabled]="!selectedFile || isUploading"
              style="margin-top: 1.5rem;"
            >
              <span *ngIf="isUploading" class="spinner"></span>
              <span *ngIf="!isUploading"><i class="fas fa-cloud-upload-alt"></i></span>
              {{ isUploading ? 'Uploading...' : 'Upload File' }}
            </button>
          </div>

          <!-- History Section -->
          <div *ngIf="activeTab === 'history'" class="card">
            <h2 style="margin-bottom: 2rem; color: var(--text-primary);">
              <i class="fas fa-history"></i> Upload History
            </h2>

            <div *ngIf="uploadedFiles.length === 0" style="text-align: center; padding: 2rem; color: var(--text-secondary);">
              <i class="fas fa-inbox" style="font-size: 3rem; margin-bottom: 1rem; display: block; opacity: 0.5;"></i>
              <p>No files uploaded yet</p>
            </div>

            <div *ngFor="let file of uploadedFiles" class="file-info" style="margin-bottom: 1rem;">
              <div class="file-icon">
                <i class="fas" [ngClass]="file.status === 'success' ? 'fa-check-circle text-success' : 'fa-clock'"></i>
              </div>
              <div class="file-details">
                <div class="file-name">{{ file.name }}</div>
                <div class="file-size">{{ file.size }} • {{ file.uploadTime }}</div>
              </div>
              <span style="color: var(--success-color); font-weight: 500;" *ngIf="file.status === 'success'">
                <i class="fas fa-check"></i> Success
              </span>
            </div>
          </div>
        </div>
      </div>
    </ng-template>
  `,
  styles: []
})
export class AppComponent {
  isLoggedIn = false;
  activeTab = 'upload';
  isLoading = false;
  isUploading = false;
  isDragging = false;
  uploadProgress = 0;

  selectedFile: File | null = null;
  alerts: Alert[] = [];
  uploadedFiles: UploadedFile[] = [];
  
  credentials = {
    email: '',
    password: ''
  };

  userEmail = '';
  userInitial = 'U';
  alertCounter = 0;

  constructor(private http: HttpClient, private cdr: ChangeDetectorRef) {
    // Check if user is already logged in (token in localStorage)
    const token = this.getToken();
    if (token) {
      this.isLoggedIn = true;
      this.userEmail = localStorage.getItem('userEmail') || '';
      this.userInitial = this.userEmail.charAt(0).toUpperCase();
    }
  }

  getToken(): string | null {
    return localStorage.getItem('jwtToken');
  }

  setToken(token: string, email: string): void {
    localStorage.setItem('jwtToken', token);
    localStorage.setItem('userEmail', email);
  }

  login() {
    console.log('Login method called');
    console.log('Credentials:', this.credentials);
    
    if (!this.credentials.email || !this.credentials.password) {
      this.showAlert('error', 'Please enter both email and password');
      return;
    }

    this.isLoading = true;
    console.log('isLoading set to true');

    console.log('Sending login request to http://localhost:8080/login');
    this.http.post<{token: string; email: string; message: string}>('http://localhost:8080/login', {
      email: this.credentials.email,
      password: this.credentials.password
    }).subscribe({
      next: (response) => {
        console.log('Login response received:', response);
        this.isLoading = false;
        console.log('isLoading set to false');
        
        this.setToken(response.token, response.email);
        console.log('Token stored');
        
        this.isLoggedIn = true;
        this.userEmail = response.email;
        this.userInitial = response.email.charAt(0).toUpperCase();
        
        console.log('State updated:', {
          isLoggedIn: this.isLoggedIn,
          userEmail: this.userEmail,
          userInitial: this.userInitial
        });
        
        this.showAlert('success', response.message || 'Login successful!');
        console.log('Alert shown, triggering change detection');
        this.cdr.detectChanges();
        console.log('Change detection triggered');
      },
      error: (error: any) => {
        console.error('Login error:', error);
        this.isLoading = false;
        const errorMessage = error.error?.error || 'Login failed. Please check your credentials.';
        this.showAlert('error', errorMessage);
        this.cdr.detectChanges();
      }
    });
  }

  logout() {
    this.isLoggedIn = false;
    this.credentials = { email: '', password: '' };
    this.selectedFile = null;
    this.alerts = [];
    this.uploadedFiles = [];
    this.activeTab = 'upload';
    localStorage.removeItem('jwtToken');
    localStorage.removeItem('userEmail');
    this.showAlert('info', 'Logged out successfully');
  }

  onDragOver(event: DragEvent) {
    event.preventDefault();
    event.stopPropagation();
    this.isDragging = true;
  }

  onDragLeave(event: DragEvent) {
    event.preventDefault();
    event.stopPropagation();
    this.isDragging = false;
  }

  onDrop(event: DragEvent) {
    event.preventDefault();
    event.stopPropagation();
    this.isDragging = false;

    const files = event.dataTransfer?.files;
    if (files && files.length > 0) {
      this.selectedFile = files[0];
      this.showAlert('info', `File selected: ${this.selectedFile.name}`);
    }
  }

  onFileSelected(event: any) {
    this.selectedFile = event.target.files[0];
    if (this.selectedFile) {
      this.showAlert('info', `File selected: ${this.selectedFile.name}`);
    }
  }

  clearFile() {
    this.selectedFile = null;
    this.uploadProgress = 0;
  }

  upload() {
    console.log('Upload method called');
    console.log('Selected file:', this.selectedFile);
    
    if (!this.selectedFile) {
      this.showAlert('error', 'Please select a file first');
      return;
    }

    this.isUploading = true;
    this.uploadProgress = 0;
    console.log('Starting upload for file:', this.selectedFile.name);

    const formData = new FormData();
    formData.append('file', this.selectedFile);

    // Simulate progress
    const progressInterval = setInterval(() => {
      if (this.uploadProgress < 90) {
        this.uploadProgress += Math.random() * 30;
      }
    }, 300);

    const token = this.getToken();
    const headers: any = {};
    if (token) {
      headers['Authorization'] = `Bearer ${token}`;
    }

    console.log('Sending upload request with token:', token ? 'yes' : 'no');

    this.http.post('http://localhost:8080/upload', formData, { headers })
      .subscribe({
        next: (res: any) => {
          console.log('Upload response received:', res);
          clearInterval(progressInterval);
          
          console.log('Setting uploadProgress to 100');
          this.uploadProgress = 100;
          
          console.log('Setting isUploading to false');
          this.isUploading = false;
          
          // Add to history immediately
          const fileName = this.selectedFile?.name || 'Unknown';
          const fileSize = this.selectedFile ? 
            (this.selectedFile.size / 1024).toFixed(2) + ' KB' : 'Unknown';
          
          console.log('Adding file to history:', fileName);
          this.uploadedFiles.unshift({
            name: fileName,
            size: fileSize,
            uploadTime: new Date().toLocaleTimeString(),
            status: 'success'
          });
          
          console.log('Showing success alert');
          this.showAlert('success', 'File uploaded successfully!');
          
          console.log('Resetting upload state');
          this.uploadProgress = 0;
          this.clearFile();
          
          console.log('Triggering change detection');
          this.cdr.detectChanges();
          
          console.log('Upload complete - isUploading:', this.isUploading);
        },
        error: (error: any) => {
          console.error('Upload error:', error);
          clearInterval(progressInterval);
          this.isUploading = false;
          this.uploadProgress = 0;
          this.cdr.detectChanges();
          
          if (error.status === 401) {
            this.showAlert('error', 'Session expired. Please login again.');
            this.logout();
          } else {
            this.showAlert('error', 'Upload failed: ' + (error.error?.error || error.message || 'Unknown error'));
          }
        }
      });
  }

  showAlert(type: 'success' | 'error' | 'warning' | 'info', message: string) {
    const id = this.alertCounter++;
    this.alerts.push({ type, message, id });

    // Auto-remove after 5 seconds
    setTimeout(() => this.removeAlert(id), 5000);
  }

  removeAlert(id: number) {
    this.alerts = this.alerts.filter(alert => alert.id !== id);
  }
}