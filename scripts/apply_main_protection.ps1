[CmdletBinding()]
param(
    [string]$Repo = "Vertex-Systems-Network/omnexa",
    [string]$Branch = "main",
    [string]$RequiredCheck = "governance",
    [ValidateRange(0, 6)]
    [int]$RequiredApprovals = 0,
    [switch]$RequireCodeOwnerReview
)

Set-StrictMode -Version Latest
$ErrorActionPreference = "Stop"

function Require-Command {
    param([Parameter(Mandatory = $true)][string]$Name)
    if (-not (Get-Command $Name -ErrorAction SilentlyContinue)) {
        throw "Required command '$Name' was not found in PATH."
    }
}

function Assert-GitHubAuthentication {
    if ($env:OMNEXA_GITHUB_ADMIN_TOKEN) {
        # gh honors GH_TOKEN without writing the credential to disk. Keep the
        # owner-controlled token short-lived and scoped to repository
        # Administration read/write for Omnexa only.
        $env:GH_TOKEN = $env:OMNEXA_GITHUB_ADMIN_TOKEN
        return
    }

    & gh auth status -h github.com *> $null
    if ($LASTEXITCODE -ne 0) {
        throw "GitHub CLI is not authenticated. Run 'gh auth login' with an admin-capable account or set OMNEXA_GITHUB_ADMIN_TOKEN for this process."
    }
}

Require-Command -Name "gh"
Assert-GitHubAuthentication

$reviewPolicy = [ordered]@{
    dismiss_stale_reviews           = $true
    require_code_owner_reviews      = [bool]$RequireCodeOwnerReview
    required_approving_review_count = $RequiredApprovals
    require_last_push_approval      = $false
}

$payload = [ordered]@{
    required_status_checks = [ordered]@{
        strict   = $true
        contexts = @($RequiredCheck)
    }
    enforce_admins                = $true
    required_pull_request_reviews = $reviewPolicy
    restrictions                  = $null
    required_conversation_resolution = $true
    allow_force_pushes            = $false
    allow_deletions               = $false
}

$json = $payload | ConvertTo-Json -Depth 10 -Compress
$endpoint = "repos/$Repo/branches/$Branch/protection"

Write-Host "Applying Omnexa branch protection to $Repo/$Branch ..."
$json | & gh api --method PUT $endpoint --input - *> $null
if ($LASTEXITCODE -ne 0) {
    throw "GitHub rejected the branch-protection update. Confirm the credential has repository Administration read/write permission."
}

Write-Host "Branch protection API update: PASS"
Write-Host "Required PR path: enabled"
Write-Host "Required check: $RequiredCheck (strict/up-to-date)"
Write-Host "Conversation resolution: required"
Write-Host "Force pushes: blocked"
Write-Host "Branch deletion: blocked"
Write-Host "Admin enforcement: enabled"
Write-Host "Required approvals: $RequiredApprovals"
Write-Host "Code Owner review required: $([bool]$RequireCodeOwnerReview)"
Write-Host ""
Write-Host "Run scripts/verify_main_protection.ps1 next and record the resulting GitHub evidence in Issue #3."
