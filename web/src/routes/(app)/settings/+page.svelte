<script lang="ts">
	import { getAuth } from '$lib/stores/auth.svelte';

	const auth = getAuth();
</script>

<div class="page">
	<div class="page-header">
		<h1>Cài đặt</h1>
		<p>Quản lý tài khoản và cài đặt trang bio</p>
	</div>

	<div class="card">
		<div class="card-header">
			<h2>Thông tin tài khoản</h2>
		</div>
		<div class="card-body">
			<div class="info-row">
				<span class="info-label">Email</span>
				<span class="info-value">{auth.user?.email}</span>
			</div>
			<div class="info-row">
				<span class="info-label">Username</span>
				<span class="info-value">@{auth.user?.username}</span>
			</div>
			<div class="info-row">
				<span class="info-label">URL trang bio</span>
				<span class="info-value url">linkbio.com/{auth.user?.username}</span>
			</div>
		</div>
	</div>

	<div class="card">
		<div class="card-header">
			<h2>Cài đặt trang</h2>
		</div>
		<div class="card-body">
			<div class="setting-row">
				<div class="setting-info">
					<span class="setting-label">Tiêu đề trang</span>
					<span class="setting-desc">Hiển thị ở đầu trang bio của bạn</span>
				</div>
				<input type="text" placeholder="Tên của bạn" class="setting-input" />
			</div>
			<div class="setting-row">
				<div class="setting-info">
					<span class="setting-label">Mô tả</span>
					<span class="setting-desc">Giới thiệu ngắn về bạn</span>
				</div>
				<textarea placeholder="Viết gì đó về bạn..." class="setting-textarea"></textarea>
			</div>
		</div>
	</div>

	<div class="card danger">
		<div class="card-header">
			<h2>Vùng nguy hiểm</h2>
		</div>
		<div class="card-body">
			<div class="danger-row">
				<div class="danger-info">
					<span class="danger-label">Xóa tài khoản</span>
					<span class="danger-desc">Xóa vĩnh viễn tài khoản và tất cả dữ liệu</span>
				</div>
				<button class="btn-danger">Xóa tài khoản</button>
			</div>
		</div>
	</div>

	<div class="coming-soon">
		<span>🚧</span>
		<p>Tính năng đang được phát triển</p>
	</div>
</div>

<style>
	.page {
		max-width: 800px;
		margin: 0 auto;
	}

	.page-header {
		margin-bottom: 1.5rem;
	}

	.page-header h1 {
		font-size: 1.25rem;
		font-weight: 600;
		color: #1f2328;
		margin-bottom: 0.25rem;
	}

	.page-header p {
		font-size: 0.875rem;
		color: #656d76;
	}

	.card {
		background: #ffffff;
		border: 1px solid #d0d7de;
		border-radius: 6px;
		margin-bottom: 1rem;
	}

	.card.danger {
		border-color: #ffcdd2;
	}

	.card.danger .card-header {
		background: #ffebee;
		border-bottom-color: #ffcdd2;
	}

	.card.danger .card-header h2 {
		color: #c62828;
	}

	.card-header {
		padding: 1rem;
		border-bottom: 1px solid #d0d7de;
	}

	.card-header h2 {
		font-size: 0.875rem;
		font-weight: 600;
		color: #1f2328;
		margin: 0;
	}

	.card-body {
		padding: 1rem;
	}

	.info-row {
		display: flex;
		justify-content: space-between;
		align-items: center;
		padding: 0.75rem 0;
		border-bottom: 1px solid #eaeef2;
	}

	.info-row:last-child {
		border-bottom: none;
	}

	.info-label {
		font-size: 0.875rem;
		color: #656d76;
	}

	.info-value {
		font-size: 0.875rem;
		color: #1f2328;
		font-weight: 500;
	}

	.info-value.url {
		font-family: ui-monospace, monospace;
		background: #f6f8fa;
		padding: 0.25rem 0.5rem;
		border-radius: 4px;
	}

	.setting-row {
		padding: 1rem 0;
		border-bottom: 1px solid #eaeef2;
	}

	.setting-row:last-child {
		border-bottom: none;
	}

	.setting-info {
		margin-bottom: 0.75rem;
	}

	.setting-label {
		display: block;
		font-size: 0.875rem;
		font-weight: 500;
		color: #1f2328;
		margin-bottom: 0.25rem;
	}

	.setting-desc {
		font-size: 0.75rem;
		color: #656d76;
	}

	.setting-input,
	.setting-textarea {
		width: 100%;
		padding: 0.5rem 0.75rem;
		border: 1px solid #d0d7de;
		border-radius: 6px;
		font-size: 0.875rem;
	}

	.setting-textarea {
		min-height: 80px;
		resize: vertical;
	}

	.setting-input:focus,
	.setting-textarea:focus {
		outline: none;
		border-color: #0969da;
		box-shadow: 0 0 0 3px rgba(9, 105, 218, 0.1);
	}

	.danger-row {
		display: flex;
		justify-content: space-between;
		align-items: center;
	}

	.danger-label {
		display: block;
		font-size: 0.875rem;
		font-weight: 500;
		color: #1f2328;
		margin-bottom: 0.25rem;
	}

	.danger-desc {
		font-size: 0.75rem;
		color: #656d76;
	}

	.btn-danger {
		background: #cf222e;
		color: white;
		padding: 0.5rem 1rem;
		border-radius: 6px;
		font-size: 0.8125rem;
		font-weight: 500;
	}

	.btn-danger:hover {
		background: #a40e26;
	}

	.coming-soon {
		text-align: center;
		padding: 2rem;
		color: #656d76;
	}

	.coming-soon span {
		font-size: 2rem;
		display: block;
		margin-bottom: 0.5rem;
	}

	.coming-soon p {
		font-size: 0.875rem;
	}
</style>
