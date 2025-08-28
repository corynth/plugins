workflow "ansible-playbook-example" {
  description = "Example using Ansible playbook with Go plugin"
}

step "run_playbook" {
  plugin = "ansible"
  action = "playbook"
  
  params = {
    playbook = <<EOF
---
- hosts: localhost
  connection: local
  gather_facts: no
  tasks:
    - name: Create directory
      file:
        path: /tmp/ansible-test
        state: directory
        
    - name: Create test file
      copy:
        content: "Hello from Ansible!"
        dest: /tmp/ansible-test/hello.txt
        
    - name: Display message
      debug:
        msg: "Ansible playbook executed successfully!"
EOF

    inventory = "localhost ansible_connection=local"
    
    vars = {
      message = "Hello from Corynth Ansible Plugin!"
      timestamp = "2024"
    }
  }
}

step "run_ad_hoc" {
  plugin = "ansible"
  action = "ad_hoc"
  depends_on = ["run_playbook"]
  
  params = {
    hosts = "localhost"
    module = "setup"
    args = "filter=ansible_os_family"
    inventory = "localhost ansible_connection=local"
  }
}