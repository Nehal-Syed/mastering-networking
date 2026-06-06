// API URL (will be proxied through Nginx)
const API_URL = 'http://localhost/api';

// Frontend is served on port 3000, backend on 8080
// We'll use Nginx as reverse proxy to avoid CORS issues

async function loadNetworkInfo() {
    try {
        const response = await fetch(`${API_URL}/network-info`);
        const data = await response.json();
        
        const networkDetails = document.getElementById('networkDetails');
        networkDetails.innerHTML = `
            <p>🔌 <strong>Your IP:</strong> ${data.client_ip || data.remote_addr}</p>
            <p>🌐 <strong>X-Forwarded-For:</strong> ${data.x_forwarded_for || 'Not set'}</p>
            <p>🏠 <strong>X-Real-IP:</strong> ${data.x_real_ip || 'Not set'}</p>
            <p>🖥️ <strong>User Agent:</strong> ${data.user_agent}</p>
            <p>📡 <strong>Protocol:</strong> ${data.protocol}</p>
            <p>🎯 <strong>Host:</strong> ${data.host}</p>
            <p>📍 <strong>Request URI:</strong> ${data.request_uri}</p>
        `;
        
        document.getElementById('networkStatus').innerHTML = '✅ Connected via Reverse Proxy';
    } catch (error) {
        console.error('Network info error:', error);
        document.getElementById('networkDetails').innerHTML = '<p>❌ Failed to load network info. Make sure backend is running.</p>';
        document.getElementById('networkStatus').innerHTML = '❌ Network Error';
    }
}

// Rest of the functions remain the same as before...
async function loadTasks() {
    try {
        const response = await fetch(`${API_URL}/tasks`);
        const tasks = await response.json();
        
        const tasksList = document.getElementById('tasksList');
        
        if (tasks.length === 0) {
            tasksList.innerHTML = '<p>No tasks yet. Create your first task!</p>';
            return;
        }
        
        tasksList.innerHTML = tasks.map(task => `
            <div class="task-card">
                <div class="task-title">${escapeHtml(task.title)}</div>
                <div class="task-desc">${escapeHtml(task.description)}</div>
                <div class="task-status ${task.completed ? 'status-completed' : 'status-pending'}">
                    ${task.completed ? '✅ Completed' : '⏳ Pending'}
                </div>
                <div class="task-actions">
                    <button onclick="toggleTask(${task.id})">
                        ${task.completed ? 'Undo' : 'Complete'}
                    </button>
                    <button class="delete-btn" onclick="deleteTask(${task.id})">
                        Delete
                    </button>
                </div>
                <small>Created: ${new Date(task.created_at).toLocaleString()}</small>
            </div>
        `).join('');
    } catch (error) {
        console.error('Load tasks error:', error);
        document.getElementById('tasksList').innerHTML = '<p>❌ Failed to load tasks. Make sure backend is running.</p>';
    }
}

async function createTask() {
    const title = document.getElementById('taskTitle').value;
    const description = document.getElementById('taskDesc').value;
    
    if (!title) {
        alert('Please enter a task title');
        return;
    }
    
    try {
        const response = await fetch(`${API_URL}/tasks`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                title: title,
                description: description,
                completed: false
            })
        });
        
        if (response.ok) {
            document.getElementById('taskTitle').value = '';
            document.getElementById('taskDesc').value = '';
            loadTasks();
        } else {
            const error = await response.text();
            alert('Failed to create task: ' + error);
        }
    } catch (error) {
        console.error('Create task error:', error);
        alert('Network error creating task');
    }
}

async function toggleTask(id) {
    try {
        const response = await fetch(`${API_URL}/tasks/${id}/toggle`, {
            method: 'PATCH'
        });
        
        if (response.ok) {
            loadTasks();
        } else {
            alert('Failed to update task');
        }
    } catch (error) {
        console.error('Toggle task error:', error);
        alert('Network error updating task');
    }
}

async function deleteTask(id) {
    if (!confirm('Are you sure you want to delete this task?')) {
        return;
    }
    
    try {
        const response = await fetch(`${API_URL}/tasks/${id}`, {
            method: 'DELETE'
        });
        
        if (response.ok) {
            loadTasks();
        } else {
            alert('Failed to delete task');
        }
    } catch (error) {
        console.error('Delete task error:', error);
        alert('Network error deleting task');
    }
}

function escapeHtml(text) {
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
}

// Initial load
loadNetworkInfo();
loadTasks();

// Refresh tasks every 10 seconds
setInterval(loadTasks, 10000);