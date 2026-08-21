[CmdletBinding()]
param(
    [string]$Repo = "Vertex-Systems-Network/omnexa",
    [string]$Branch = "main",
    [string]$RequiredCheck = "governance",
    [ValidateRange(0, 6)]
    [int]$MinimumApprovals = 0,
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
        $env:GH_TOKEN = $env:OMNEXA_GITHUB_ADMIN_TOKEN
        return
    }

    & gh auth status -h github.com *> $null
    if ($LASTEXITCODE -ne 0) {
        throw "GitHub CLI is not authenticated. Run 'gh auth login' or set OMNEXA_GITHUB_ADMIN_TOKEN for this process."
    }
}

function Assert-True {
    param(
        [Parameter(Mandatory = $true)][bool]$Condition,
        [Parameter(Mandatory = $true)][string]$Message
    )
    if (-not $Condition) {
        throw "FAIL: $Message"
    }
    Write-Host "PASS: $Message"
}

Require-Command -Name "gh"
Assert-GitHubAuthentication

$branchJson = & gh api "repos/$Repo/branches/$Branch"
if ($LASTEXITCODE -ne 0) {
    throw "Unable to read branch metadata for $Repo/$Branch."
}
$branchData = $branchJson | ConvertFrom-Json
Assert-True -Condition ([bool]$branchData.protected) -Message "branch reports protected=true"

$protectionJson = & gh api "repos/$Repo/branches/$Branch/protection"
if ($LASTEXITCODE -ne 0) {
    throw "Unable to read branch protection for $Repo/$Branch."
}
$protection = $protectionJson | ConvertFrom-Json

Assert-True -Condition ($null -ne $protection.required_status_checks) -Message "required status checks are configured"
Assert-True -Condition ([bool]$protection.required_status_checks.strict) -Message "required status checks are strict/up-to-date"
$contexts = @($protection.required_status_checks.contexts)
Assert-True -Condition ($contexts -contains $RequiredCheck) -Message "required status check '$RequiredCheck' is configured"

Assert-True -Condition ($null -ne $protection.required_pull_request_reviews) -Message "pull-request review protection is configured"
$approvalCount = [int]$protection.required_pull_request_reviews.required_approving_review_count
Assert-True -Condition ($approvalCount -ge $MinimumApprovals) -Message "required approvals >= $MinimumApprovals (actual: $approvalCount)"
Assert-True -Condition ([bool]$protection.required_pull_request_reviews.dismiss_stale_reviews) -Message "stale approvals are dismissed"

if ($RequireCodeOwnerReview) {
    Assert-True -Condition ([bool]$protection.required_pull_request_reviews.require_code_owner_reviews) -Message "Code Owner review is required"
}

Assert-True -Condition ([bool]$protection.enforce_admins.enabled) -Message "branch protection is enforced for administrators"
Assert-True -Condition ([bool]$protection.required_conversation_resolution.enabled) -Message "conversation resolution is required"
Assert-True -Condition (-not [bool]$protection.allow_force_pushes.enabled) -Message "force pushes are blocked"
Assert-True -Condition (-not [bool]$protection.allow_deletions.enabled) -Message "branch deletion is blocked"

Write-Host ""
Write-Host "Omnexa main branch protection verification: PASS"
Write-Host "Repository: $Repo"
Write-Host "Branch: $Branch"
Write-Host "Required check: $RequiredCheck"
Write-Host "Required approvals: $approvalCount"
Write-Host ""
Write-Host "Configuration verification does not replace the controlled negative tests tracked in Issue #3 (direct push, force push, deletion, failed governance PR and CODEOWNERS behavior)."
