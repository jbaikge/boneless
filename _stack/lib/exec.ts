import * as childProcess from 'child_process';

const exec = function(command: string, options?: childProcess.SpawnSyncOptions) {
  const proc = childProcess.spawnSync('bash', ['-c', command], options);

  if (proc.error) {
    throw proc.error;
  }

  if (proc.status != 0) {
    if (proc.stdout || proc.stderr) {
      throw new Error(`[Status ${proc.status}] stdout: ${proc.stdout?.toString().trim()}\n\n\nstderr: ${proc.stderr?.toString().trim()}`);
    }
    throw new Error(`process exited with status ${proc.status}`);
  }

  return proc;
}

export default exec;